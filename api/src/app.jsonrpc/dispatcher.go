package app_jsonrpc

import (
	base "app.base"
	models "app.models"
	"context"
	"encoding"
	"encoding/gob"
	"fmt"
	"reflect"
	"strings"
)

var (
	gobEncoderType        = reflect.TypeOf((*gob.GobEncoder)(nil)).Elem()
	gobDecoderType        = reflect.TypeOf((*gob.GobDecoder)(nil)).Elem()
	binaryMarshalerType   = reflect.TypeOf((*encoding.BinaryMarshaler)(nil)).Elem()
	binaryUnmarshalerType = reflect.TypeOf((*encoding.BinaryUnmarshaler)(nil)).Elem()

	errt = reflect.TypeOf((*error)(nil)).Elem()

	validatedTypes []*validatedType
)

type validatedType struct {
	t   reflect.Type
	err *error
}

type funcData struct {
	inNum int
	reqt  reflect.Type
	fv    reflect.Value
}

type serviceData struct {
	sv      reflect.Value
	funcMap map[string]*funcData
}

type HandlerFunc func(clientAddr string, request interface{}) (response interface{})
type JsonHandlerFunc func(clientAddr string, request *models.RpcJsonRequest) (response *models.RpcJsonResponse)

type dispatcherRpc struct {
	ctx    context.Context
	cancel context.CancelFunc

	serviceMap   map[string]*serviceData
	methodPrefix string
	handler      JsonHandlerFunc

	log *base.Logger
}

func newDispatcher(ctx context.Context, prefix string) *dispatcherRpc {
	disp := &dispatcherRpc{
		serviceMap:   make(map[string]*serviceData),
		methodPrefix: prefix,
	}

	disp.ctx, disp.cancel = context.WithCancel(ctx)
	disp.log = base.GetLog(disp.ctx, disp)

	return disp
}

func (d *dispatcherRpc) addService(serviceName string, service interface{}) error {
	var err error
	if serviceName == "" {
		return fmt.Errorf("dispatcherRpc.AddService: serviceName cannot be empty")
	}

	if _, ok := d.serviceMap[serviceName]; ok {
		return fmt.Errorf("dispatcherRpc.AddService service with name=[%s] has been already registered", serviceName)
	}

	funcMap := make(map[string]*funcData)

	st := reflect.TypeOf(service)
	if st.Kind() == reflect.Struct {
		return fmt.Errorf("dispatcherRpc.AddService service [%s] must be a pointer to struct, i.e. *%s", serviceName, st)
	}
	for i := 0; i < st.NumMethod(); i++ {
		mv := st.Method(i)

		if mv.PkgPath != "" {
			continue
		}

		if d.methodPrefix != "" && !strings.HasPrefix(mv.Name, d.methodPrefix) {
			continue
		}

		funcName := serviceName + "." + mv.Name
		fd := &funcData{
			fv: mv.Func,
		}
		if fd.inNum, fd.reqt, err = validateFunc(funcName, fd.fv, true); err != nil {
			return fmt.Errorf("dispatcherRpc.AddService validate func [%s] error: %s", funcName, err)
		}
		funcMap[mv.Name] = fd
	}
	if len(funcMap) == 0 {
		return fmt.Errorf("dispatcherRpc.AddService the service %s has no methods suitable for rpc", serviceName)
	}

	d.serviceMap[serviceName] = &serviceData{
		sv:      reflect.ValueOf(service),
		funcMap: funcMap,
	}

	return nil
}

func (d *dispatcherRpc) newJsonHandlerFunc() (JsonHandlerFunc, error) {
	if len(d.serviceMap) == 0 {
		return nil, fmt.Errorf("Dispatcher.NewHandlerFunc register at least one service before calling HandlerFunc()")
	}

	serviceMap := copyServiceMap(d.serviceMap)
	return func(clientAddr string, request *models.RpcJsonRequest) *models.RpcJsonResponse {
		return dispatchJsonRequest(d.ctx, serviceMap, request)
	}, nil
}

func dispatchJsonRequest(ctx context.Context, serviceMap map[string]*serviceData, req *models.RpcJsonRequest) *models.RpcJsonResponse {
	log := ctx.Value("logger").(*base.Logger)
	names := strings.Split(req.Method, ".")
	if len(names) != 2 {
		log.Error(fmt.Sprintf("dispatchJsonRequest mistake requst method [%s]", req.Method))
		return models.PrepareRpcErrorResponse(req.Id, IncorrectRequestTypeError, req.UuidWsClients)
	}
	serviceName, funcName := names[0], fmt.Sprintf("%s%s", DefaultMethodPrefix, names[1])

	s, ok := serviceMap[serviceName]

	if !ok {
		log.Error(fmt.Sprintf("dispatchJsonRequest unknown service name [%s]", serviceName))
		return models.PrepareRpcErrorResponse(req.Id, IncorrectRequestTypeError, req.UuidWsClients)
	}

	fd, ok := s.funcMap[funcName]
	if !ok {
		log.Error(fmt.Sprintf("dispatchJsonRequest unknown method [%s.%s]", serviceName, funcName))
		return models.PrepareRpcErrorResponse(req.Id, UnknownMethodError, req.UuidWsClients)
	}
	var inArgs []reflect.Value
	if fd.inNum > 0 {
		inArgs = make([]reflect.Value, fd.inNum)
		dt := 0
		if serviceName != "" {
			dt = 1
			inArgs[0] = s.sv
		}
		if fd.inNum > dt {
			reqv := reflect.ValueOf(req.Params)
			reqt := reflect.TypeOf(req.Params)
			if reqt != fd.reqt {
				log.Error(fmt.Sprintf("dispatchJsonRequest unexpected request type for method [%s]: %s. Expected %s", funcName, reqt, fd.reqt))
				return models.PrepareRpcErrorResponse(req.Id, UnexpectedRequestTypeError, req.UuidWsClients)
			}
			inArgs[len(inArgs)-1] = reqv
		}
	}

	outArgs := fd.fv.Call(inArgs)

	resp := &models.RpcJsonResponse{
		Id:            req.Id,
		UuidWsClients: req.UuidWsClients,
	}

	if len(outArgs) == 1 {
		if isErrorType(outArgs[0].Type()) {
			if !outArgs[0].IsNil() {
				resp := models.PrepareRpcErrorResponse(req.Id, IncorrectRequestTypeError, req.UuidWsClients)
				resp.Error = &models.RpcJsonError{
					Code:    500,
					Message: getErrorString(outArgs[0]),
				}
			}
		} else {
			resp.Result, ok = outArgs[0].Interface().([]byte)
			if !ok {
				log.Error(fmt.Sprintf("dispatchJsonRequest can't convert response from [%s]: %s.", funcName, outArgs[0].Interface()))
			}
		}
	} else if len(outArgs) == 2 {
		resp.Result, ok = outArgs[0].Interface().([]byte)
		if !ok {
			log.Error(fmt.Sprintf("dispatchJsonRequest can't convert response from [%s]: %s.", funcName, outArgs[0].Interface()))
		}

		if !outArgs[1].IsNil() {
			resp.Error = &models.RpcJsonError{
				Code:    500,
				Message: getErrorString(outArgs[1]),
			}
		}
	}

	return resp
}

func getErrorString(v reflect.Value) string {
	if v.IsNil() {
		return ""
	}
	return v.Interface().(error).Error()
}

func RegisterType(x interface{}) {
	gob.Register(x)
}

func registerType(s, funcName string, t reflect.Type) error {
	if t.Kind() == reflect.Struct {
		return fmt.Errorf("%s in the function [%s] should be passed by reference, i.e. *%s", s, funcName, t)
	}
	if err := validateType(t); err != nil {
		return fmt.Errorf("%s in the function [%s] cannot contain %s", s, funcName, err)
	}

	t = removePtr(t)
	tv := reflect.New(t)
	if t.Kind() != reflect.Struct {
		tv = reflect.Indirect(tv)
	}

	switch t.Kind() {
	case reflect.Array, reflect.Slice, reflect.Map, reflect.Struct:
		RegisterType(tv.Interface())
	default:
	}

	return nil
}

func validateFunc(funcName string, fv reflect.Value, isMethod bool) (inNum int, reqt reflect.Type, err error) {
	if funcName == "" {
		err = fmt.Errorf("func name can't be empty")
		return
	}

	ft := fv.Type()
	if ft.Kind() != reflect.Func {
		err = fmt.Errorf("function [%s] must be a function instead of %s", funcName, ft)
		return
	}

	inNum = ft.NumIn()
	outNum := ft.NumOut()

	fmt.Println("inNum: ", inNum)
	fmt.Println("outNum: ", outNum)

	dt := 0
	if isMethod {
		dt = 1
	}

	if inNum == 2+dt {
		fmt.Println(" 2+dt ft.In(dt).Kind(): ", ft.In(dt).Kind())
		if ft.In(dt).Kind() != reflect.String {
			err = fmt.Errorf("unexpected type for the first argument of the function [%s]: [%s]. Expected string", funcName, ft.In(dt))
			return
		}
	} else if inNum > 2+dt {
		err = fmt.Errorf("unexpected number of arguments in the function [%s]: %d. Expected 0, 1 (request) or 2 (clientAddr, request)", funcName, inNum-dt)
		return
	}

	if outNum == 2 {
		if !isErrorType(ft.Out(1)) {
			err = fmt.Errorf("unexpected type for the second return value of the function [%s]: [%s]. Expected [%s]", funcName, ft.Out(1), errt)
			return
		}
	} else if outNum > 2 {
		err = fmt.Errorf("unexpected number of return values for the function %s: %d. Expected 0, 1 (response) or 2 (response, error)", funcName, outNum)
		return
	}

	if inNum > dt {
		reqt = ft.In(inNum - 1)
		if err = registerType("request", funcName, reqt); err != nil {
			return
		}
	}

	if outNum > 0 {
		respt := ft.Out(0)
		if !isErrorType(respt) {
			if err = registerType("response", funcName, ft.Out(0)); err != nil {
				return
			}
		}
	}

	return
}

func validateType(t reflect.Type) (err error) {
	t = removePtr(t)
	for _, vd := range validatedTypes {
		if vd.t == t {
			return *vd.err
		}
	}
	validatedTypes = append(validatedTypes, &validatedType{
		t:   t,
		err: &err,
	})

	switch t.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.UnsafePointer:
		err = fmt.Errorf("%s. Found [%s]", t.Kind(), t)
		return err
	case reflect.Array, reflect.Slice:
		if err = validateType(t.Elem()); err != nil {
			err = fmt.Errorf("%s in the %s [%s]", err, t.Kind(), t)
			return err
		}
	case reflect.Map:
		if err = validateType(t.Elem()); err != nil {
			err = fmt.Errorf("%s in the value of map [%s]", err, t)
			return err
		}
		if err = validateType(t.Key()); err != nil {
			err = fmt.Errorf("%s in the key of map [%s]", err, t)
			return err
		}
	case reflect.Struct:
		if supportsGob(t) {
			return nil
		}
		// Special case for struct{}
		if t.NumField() == 0 {
			return nil
		}
		n := 0
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			if f.PkgPath == "" {
				if err = validateType(f.Type); err != nil {
					err = fmt.Errorf("%s in the field [%s] of struct [%s]", err, f.Name, t)
					return err
				}
				n++
			}
		}
		if n == 0 {
			err = fmt.Errorf("struct without exported fields [%s]", t)
			return err
		}
	}

	return err
}

func copyServiceMap(sm map[string]*serviceData) map[string]*serviceData {
	serviceMap := make(map[string]*serviceData)
	for sk, sv := range sm {
		funcMap := make(map[string]*funcData)
		for fk, fv := range sv.funcMap {
			funcMap[fk] = fv
		}
		serviceMap[sk] = &serviceData{
			sv:      sv.sv,
			funcMap: funcMap,
		}
	}
	return serviceMap
}

func supportsGob(t reflect.Type) bool {
	if t.Kind() != reflect.Struct {
		panic(fmt.Sprintf("non-struct type passed to supportsGob: %s", t))
	}
	t = reflect.PtrTo(t)
	if t.Implements(gobEncoderType) && t.Implements(gobDecoderType) {
		return true
	}
	if t.Implements(binaryMarshalerType) && t.Implements(binaryUnmarshalerType) {
		return true
	}
	return false
}

func isErrorType(t reflect.Type) bool {
	return t.Implements(errt)
}

func removePtr(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

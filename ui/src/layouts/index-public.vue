<template>
    <ion-page>
        <ion-router-outlet />
    </ion-page>
</template>

<script>
import { IonPage, IonRouterOutlet } from '@ionic/vue';

export default {
    name: 'LayoutsIndexPublic',
    components: { IonPage, IonRouterOutlet }
}
</script>

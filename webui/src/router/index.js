import {createRouter, createWebHashHistory} from 'vue-router'
import login from '../views/login.vue'
import ProfileComponent from '../views/ProfileComponent.vue'
import HomeView from '../views/HomeView.vue'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/', component: login},
		{
			path: '/profile/:username',
			
			component: ProfileComponent
			
		},
		
		{path: '/link1', component: HomeView},
		{path: '/link2', component: HomeView},
		{path: '/some/:id/link', component: HomeView},
	]
})

export default router

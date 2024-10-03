import {createRouter, createWebHashHistory} from 'vue-router'
import login from '../views/login.vue'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/', component: login},
		{
			path: '/profile/:username',
			name: 'Profile',
			component: ProfileComponent,
			props: true
		},
		
		{path: '/link1', component: HomeView},
		{path: '/link2', component: HomeView},
		{path: '/some/:id/link', component: HomeView},
	]
})

export default router

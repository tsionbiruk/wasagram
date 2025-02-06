import {createRouter, createWebHashHistory} from 'vue-router'
import HomeView from '../views/HomeView.vue'
import LoginView from '../views/LoginView.vue'
import ProfileView from '../views/ProfileView.vue'
import PhotoView from '../views/PhotoView.vue'
import UsersView from '../views/UsersView.vue'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/',               component: HomeView},
		{path: '/login',          component: LoginView},
		{path: '/profile',        component: ProfileView},
		{path: '/view',           component: PhotoView},
		{path: '/users',          component: UsersView},
	]
})

export default router

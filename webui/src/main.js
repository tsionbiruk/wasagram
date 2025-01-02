import {createApp, ref} from 'vue'
import App from './App.vue'
import router from './router'
import axios from './services/axios.js';
import ErrorMsg from './components/ErrorMsg.vue';
import LoadingSpinner from './components/LoadingSpinner.vue';
import StreamPost from './components/StreamPost.vue';
import LikeButton from './components/LikeButton.vue';
import CommentSection from './components/CommentSection.vue';
import Comment from './components/Comment.vue';
import UserOptions from './components/UserOptions.vue';

import './assets/dashboard.css'
import './assets/main.css'

const app = createApp(App);
app.config.globalProperties.$token = ref(0);
app.config.globalProperties.$username = ref(null);
app.config.globalProperties.$photo = ref(null);
app.config.globalProperties.$photoorigin = ref(null);
app.config.globalProperties.$axios = axios;

app.config.globalProperties.$timestamp2date = timestamp => {
	var d    = new Date(timestamp * 1000),
			yyyy = d.getFullYear(),
			mm   = ('0' + (d.getMonth() + 1)).slice(-2),
			dd   = ('0' + d.getDate()).slice(-2),
			hh   = d.getHours(),
			min  = ('0' + d.getMinutes()).slice(-2),
			ampm = 'AM',
			time;
	time = dd + '.' + mm + '.' + yyyy + ', ' + hh + ':' + min;
	return time;
};

app.component("ErrorMsg", ErrorMsg);
app.component("LoadingSpinner", LoadingSpinner);
app.component("StreamPost", StreamPost);
app.component("LikeButton", LikeButton);
app.component("CommentSection", CommentSection);
app.component("Comment", Comment);
app.component("UserOptions", UserOptions);
app.use(router);
app.mount('#app');

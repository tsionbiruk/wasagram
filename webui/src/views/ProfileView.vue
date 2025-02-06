<script setup>
</script>
<script>
export default {
	data: function() {
		return {
			photodata: [],
			following: [],
			followers: [],
			banned:    [],
			errormsg: null,
			loading: false,
			form: {
				username: null,
			},
		}
	},
	methods: {
		async refresh() {
			this.loading = true;
			this.errormsg = null;
			this.photodata = [];
			this.$axios({
				method: 'get',
				url: `/users/${this.$username.value}/profile`,
				headers: {Authorization: this.$token.value},
			})
			.then(async response => {
				this.errormsg = null;
				let photos = response.data.photos;
				this.following = response.data.following;
				this.followers = response.data.followers;
				this.banned    = response.data.banned;
				for (let i = 0; i < photos.length; i++) {
					let id         = photos[i].id,
					    author     = photos[i].user,
					    ts         = photos[i].timestamp,
					    comments   = photos[i].comments,
					    likes      = photos[i].likes,
					    png64      = '',
					    response   = null;
					try {
						response = await this.$axios.get(`/photos/${id}`, {responseType: 'arraybuffer',});
					} catch(error) {
						this.errormsg = `Failed to retrieve photo ${id}: ${error.response.data}`;
						console.error(error);
					}
					if (response.data) {
						let binary = '';
						let bytes = new Uint8Array(response.data);
						let len = bytes.byteLength;
						for (let i = 0; i < len; i++) {
							binary += String.fromCharCode(bytes[i]);
						}
						png64 = btoa(binary);
					}
					this.photodata.push({
						id: id,
						png64: png64,
						author: author,
						timestamp: ts,
						comments: comments,
						likes: likes,
					});
				}
			})
			.catch(error => {
				this.errormsg = `Failed to retrieve profile data: ${error.response.data}`
				console.error(error);
				this.loading = false;
				return
			});
			this.loading = false;
		},
		rename() {
			this.$axios({
				method: 'post',
				url: `/users/${this.$username.value}/username`,
				data: JSON.stringify(this.form.username),
				headers: {Authorization: this.$token.value},
			})
			.then(_ => {
				this.errormsg = null;
				this.$username.value = this.form.username;
				this.refresh();
			})
			.catch(error => {
				this.errormsg = `User rename failed: ${error.response.data}`;
				console.error(error.response);
			});
		},
		upload() {
			let file = this.$refs.png.files[0];
			const reader = new FileReader();
			reader.onerror = error => {
				this.errormsg = `Failed to read photo file: ${error.response.data}`;
				console.error(error.response);
			}
			reader.onload = res => {
				this.$axios({
					method: 'post',
					url: `/users/${this.$username.value}/photos`,
					data: res.target.result,
					headers: {Authorization: this.$token.value},
				})
				.then(_ => {
					this.errormsg = null;
					this.refresh();
				})
				.catch(error => {
					this.errormsg = `Photo upload failed: ${error.response.data}`;
					console.error(error.response);
				});
			};
			reader.readAsArrayBuffer(file);
		}
	},
	mounted() {
		this.$photoorigin.value = 'profile';
		this.refresh()
	}
}
</script>

<template>
	<div>
		<div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
		
			<h1 class="h2">My Profile: {{ this.$username.value }}</h1>
			<form v-on:submit="rename">
					<label for="form-username">Change Username</label>
					<br>
					<input v-model="form.username" id="form-username">
					<button class="btn btn-sm btn-outline-secondary" type="submit">Confirm</button>
			</form>	
            			
			<div>
				<button type="button" class="btn btn-sm btn-outline-secondary" @click="refresh">
					Refresh
				</button>
				
			</div>
		</div>
        <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
		<div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-top pt-3 pb-2 mb-3 border-bottom" style="margin-bottom: 20px">
			<figure>
					<figcaption><b>Following({{ this.following.length }}):</b></figcaption>
					<ul>
						<li v-for="name in this.following" :key="name">{{ name }}</li>
					</ul>
				</figure>
				<figure>
					<figcaption><b>Followers ({{ this.followers.length }}):</b></figcaption>
					<ul >
						<li v-for="name in this.followers" :key="name">{{ name }}</li>
					</ul>
				</figure>
				<figure>
					<figcaption><b>Banned  ({{ this.banned.length }}):</b></figcaption>
					<ul>
						<li v-for="name in this.banned" :key="name">{{ name }}</li>
					</ul>
				</figure>
			</div>
			<h5>My photo stream:</h5>
			<div class="d-flex flex-column flex-wrap flex-md-nowrap align-items-between justify-items-center" style="max-width: 1000px; background-color: rgb(212, 204, 188) ; border-radius: 10px; margin: 0 auto; margin-top: 20px; padding: 20px;">
				<h5>Upload photo:</h5>
				<input ref="png" @change="upload" type="file">
			</div>
			
			<div v-if="this.photodata.length">
				<div v-for="photo in this.photodata" :key="-photo.timestamp">
					<StreamPost :params="photo" @deleted="refresh"></StreamPost>
				</div>
			</div>
			<div v-else-if="this.$token.value">
				<p><br>Your own photo stream is empty!</p>
			</div>
		</div>
</template>

<style>
</style>

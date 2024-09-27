<script>
export default {
	props: ['profile'],
	data: function() {
		return {
			photodata: [],
			errormsg: null,
			loading: false,
		}
	},
	methods: {
		async refresh() {
			this.loading = true;
			this.errormsg = null;
			this.photodata = [];
			
			try {
				const response = await this.$axios({
					method: 'get',
					url: `/users/${this.$username.value}/stream/${this.profile.username}`,
					headers: { Authorization: this.$token.value },
				});
				
				let photos = response.data;
				for (let i = 0; i < photos.length; i++) {
					let id         = photos[i].PhotoId,
						caption	   = photos[i].Caption,
					   
					    ts         = photos[i].timestamp,
					    comments   = photos[i].comments,
						comment_count = photos[i].comment_count,
					    likes      = photos[i].likes,
						like_count = photos[i].like_count,
					    png64      = '',
					    response   = null;

					try {
						const photoResponse = await this.$axios.get(`/photos/${id}`, { responseType: 'arraybuffer' });
						let binary = '';
						let bytes = new Uint8Array(photoResponse.data);
						for (let j = 0; j < bytes.byteLength; j++) {
							binary += String.fromCharCode(bytes[j]);
						}
						png64 = btoa(binary); // Convert ArrayBuffer to base64
					} catch (error) {
						this.errormsg = `Failed to retrieve photo ${id}: ${error.response.data}`;
						console.error(error);
					}

					// Push photo data to array
					this.photodata.push({
						id: id,
						png64: png64,
						caption: caption,
						
						timestamp: ts,
						comments: comments,
						comment_count: comment_count,
						likes: likes,
						like_count : like_count ,
					});
				}
			} catch (error) {
				this.errormsg = `Failed to retrieve photo stream: ${error.response.data}`;
				console.error(error);
			} finally {
				this.loading = false; // Always set loading to false after request
			}
		},
	},
	mounted() {
		this.$photoorigin.value = 'stream';
		this.refresh();
	},
}
</script>


<template>
	<div>
		<div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
			<h1 class="h2">Home Page</h1>
			<div v-if="this.$token.value" class="btn-toolbar mb-2 mb-md-0">
				<div class="btn-group me-2">
					<button type="button" class="btn btn-sm btn-outline-secondary" @click="refresh">
						Refresh
					</button>
				</div>
			</div>
		</div>

		<div v-if="this.$token.value && this.photodata.length" style = " align-items: center;">
			<div v-for="photo in this.photodata" :key="-photo.timestamp">
				<StreamPost :params="photo" @deleted="refresh"></StreamPost>
			</div>
		</div>
		<div v-else-if="this.$token.value">
			<p>Your feed is empty.</p>
		</div>
		<div v-else>
			<p>Please <RouterLink to="/login">log in</RouterLink> to access your feed.</p>
		</div>

		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
	</div>
</template>

<style>
</style>
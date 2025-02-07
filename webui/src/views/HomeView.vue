<script>
export default {
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
			this.$axios({
				method: 'get',
				url: `/users/${this.$username.value}/stream`,
				headers: {Authorization: this.$token.value},
			})
			.then(async response => {
				this.errormsg = null;
				let photos = response.data;
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
						// Convert ArrayBuffer to base64, credit: https://stackoverflow.com/a/9458996/10792539
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
				this.errormsg = `Failed to retrieve photo stream: ${error.response.data}`
				console.error(error);
				this.loading = false;
				return
			});
			this.loading = false;
		},
	},
	mounted() {
		this.$photoorigin.value = 'stream';
		this.refresh();
	}
}
</script>

<template>
	<div>
		<div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
			<h2 class="h2">Home Page</h2>
			<div v-if="this.$token.value" class="btn-toolbar mb-2 mb-md-0">
				<div class="btn-group me-2">
					<button type="button" class="btn btn-sm btn-outline-secondary" @click="refresh">
						Refresh
					</button>
				</div>
			</div>
		</div>

		<div v-if="this.$token.value && this.photodata.length" class="photo-grid">
			<div v-for="photo in this.photodata" :key="-photo.timestamp" class="photo-item">
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
.photo-grid {
	display: grid;
	grid-template-columns: repeat(auto-fill, minmax(700px, 1fr)); /* Dynamic column count */
	gap: 10px; /* Space between items */
	justify-content: center;
	padding: 10px;
}

.photo-item {
	overflow: hidden;
	border-radius: 10px;
}
</style>


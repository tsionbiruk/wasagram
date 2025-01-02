<script>
export default {
	data: function() {
		return {
			errormsg: null,
			loading: false,
		}
	},
	created() {
		this.refresh();
	},
	methods: {
		back() {
			this.$router.go(-1);
		},
		async refresh() {
			this.loading = true;
			this.errormsg = null;
			let id = this.$photo.value.id;
			await this.$axios({
				method: 'get',
				url: `/users/${this.$username.value}/${this.$photoorigin.value}`,
				headers: {Authorization: this.$token.value},
			})
			.then(async response => {
				this.errormsg = null;
				let photos = this.$photoorigin.value == 'stream' ? response.data : response.data.photos;
				for (let i = 0; i < photos.length; i++) {
					if (photos[i].id != id) {
						continue;
					}
					let author     = photos[i].user,
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

					this.$photo.value = {
						id: id,
						png64: png64,
						author: author,
						timestamp: ts,
						comments: comments,
						likes: likes,
					};
				}
			})
			.catch(error => {
				this.errormsg = `Failed to retrieve ${this.$photoorigin.value}: ${error.response.data}`
				console.error(error);
			});
			this.loading = false;
		},
	},
}
</script>

<template>
	<div style="justify-items: center">
		<StreamPost :params="this.$photo.value" @deleted="back"></StreamPost>
		<div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
			<LikeButton v-if="this.$photo.value.author != this.$username.value" @changed="refresh"></LikeButton>	
			<button v-if="this.$token.value" type="button" class="btn btn-sm btn-outline-secondary" @click="back">
						Back
					</button>
			
		</div>	
		<CommentSection @changed="refresh"></CommentSection>
		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
	</div>
</template>

<style>
</style>

<script>
export default {
	props: ['params'],
	data: function() {
		return {
			datestr: this.$timestamp2date(this.params.timestamp)
		}
	},
	methods: {
		view() {
			this.$photo.value = this.params;
			this.$router.push('/view');
		},
		remove() {
			this.$axios({
				method: 'delete',
				url: `/users/${this.$username.value}/photos/${this.params.id}`,
				headers: {Authorization: this.$token.value},
			})
			.then(_ => {
				this.$photo.value = null;
				this.$emit('deleted');
			})
			.catch(error => {
				console.error(error.response);
			});
		}
	},
}
</script>

<template>
<div class="d-flex flex-column flex-wrap flex-md-nowrap align-items-between justify-items-center" style="max-width: 1000px; background-color: rgb(212, 204, 188) ; border-radius: 10px; margin: 0 auto; margin-top: 20px; padding: 20px;">
	<div class="align-self-center">
		<img class="photo" :src="`data:image/png;base64,${params.png64}`" :alt="`photo${params.photoid}`" @click="this.view()">
	</div>	
	<div class="d-flex justify-content-between flex-wrap flex-md-nowrap" style="padding:20px;">
		<div style="float:left">
			<p><b>{{ this.datestr }}</b></p>
			<p><b>Author:</b> {{ params.author }}</p>
		</div>
		<div style="float:right">
			<p><b>Comments:</b> {{ params.comments.length }}</p>
			<p><b>Likes:</b> {{ params.likes.length }}</p>
			<button v-if="params.author == this.$username.value" type="button" class="btn btn-sm btn-danger" @click="remove">
				Delete
			</button>
		</div>
	</div>
</div>
</template>

<style>

.photo {
	min-height: 20vh;
	max-height: 50%;
	max-width: 50vh;
	width: auto;
	height: auto;
	border-radius: 5px;
}

</style>

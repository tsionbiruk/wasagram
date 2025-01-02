<script>
export default {
	data: function() {
		return {
			form: {
				text: null,
			}
		}
	},
	methods: {
		post() {
			this.$axios({
				method: 'post',
				url: `/users/${this.$username.value}/comments/${this.$photo.value.id}`,
				data: JSON.stringify(this.form.text),
				headers: {Authorization: this.$token.value},
			})
			.then(_ => {
				this.form.text = null;
				this.$emit('changed');
			})
			.catch(error => {
				console.error(error.response);
			});
	 }
	},
}
</script>

<template>
<div style = "float: center">
	<h5>Comments:</h5>
	<div class = "publish">
	<form @submit="post">
		<input v-model="form.text" id="form-text"/>
		<button style ="float:right" class="btn btn-sm btn-primary" type="submit">Publish comment</button>
	</form>
	</div>
	<div v-for="comment in this.$photo.value.comments" :key="-comment.timestamp">
		<Comment :params="comment" @changed="this.$emit('changed')"></Comment>
	</div>
</div>
</template>

<style>
.publish {
	justify-items: center;
	align-items: center;
	margin: 0 auto;
	margin-top: 20px;
	padding: 30px;
	background-color: rgb(255, 243, 220);
	max-width: 800px;
	border-radius: 10px;
}

</style>

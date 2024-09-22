<script>
import deleteComments from 'webui/src/components/delete_comments.vue'; 
export default {
	components: {
    deleteComments,  // Register PostComments component
    
  },
	
  	props: ['params'],
	data: function() {
		return {
			form: {
				text: null, 
			},
			isLoading: false, 
			errorMessage: null 
		}
	},
	methods: {
		
		validateForm() {
			// Check if the comment text is empty
			if (!this.form.text || this.form.text.trim() === '') {
				this.errorMessage = 'Comment cannot be empty';
				return false;
			}
			this.errorMessage = null; 
			return true;
		},
		
	
		post() {
			
			if (!this.validateForm()) {
				return;
			}

			// Set the loading state to true to indicate the request is being processed
			this.isLoading = true;
			
			// Make the POST request
			this.$axios({
				method: 'post',
				url: `/users/${this.$username.value}/comments/${this.$photo.value.PhotoId}`,
				data: JSON.stringify(this.form.text),
				headers: { Authorization: this.$token.value },
			})
			.then(_ => {
				this.form.text = null; 
				this.$emit('changed'); 
			})
			.catch(error => {
				console.error(error.response); 
			})
			.finally(() => {
				this.isLoading = false; 
			});
		}
	},
}
</script>


<template>
<div style="max-width: 800px; margin: 0 auto;">
	<!-- Comments section header -->
	<h5>Comments:</h5>
		<div class="publish">
			
			<form @submit.prevent="post">
				<input v-model="form.text" id="form-text" placeholder="Write a comment..." class="form-control" />
				<!-- Add some margin and align the button to the right -->
				<button style="float:right; margin-top: 10px;" class="btn btn-sm btn-primary" type="submit">Publish comment</button>
			</form>
			<p>Posting a comment for photo: {{ params.PhotoId }}</p>
		</div>
	
		<!-- Iterate through each comment and pass it to the Comment component -->
		<div v-for="comment in this.$photo.value.Comments " :key="comment.Id">
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
	background-color: rgb(218, 202, 173);
	max-width: 800px;
	border-radius: 10px;
}

</style>
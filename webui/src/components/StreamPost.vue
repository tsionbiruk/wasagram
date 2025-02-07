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
<div class="d-flex flex-column flex-wrap flex-md-nowrap align-items-between justify-items-center photo-card">
	<div class="photo-container">
		<img class="photo" :src="`data:image/png;base64,${params.png64}`" :alt="`photo${params.photoid}`" @click="this.view()">
	</div>	
	<div class="d-flex justify-content-between flex-wrap flex-md-nowrap photo-details">
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


.photo-card {
  max-width: 700px;
  background-color: rgb(159, 182, 225);
  border-radius: 25px;
  margin: 20px auto;
  padding: 10px;
  display: flex;
  flex-direction: column;
  align-items: center;
}


.photo-container {
  padding: 20px;
  border: 2px solid #fff; 
  border-radius: 20px;
  box-shadow: 0px 4px 8px rgba(0, 0, 0, 0.2); 
  background: #ffffff; 
  display: flex;
  justify-content: center;
  align-items: center;
}


.photo {
  max-width: 100%;
  max-height: 300px; 
  border-radius: 20px;
  cursor: pointer;
}


.photo-details {
  width: 100%;
  padding: 10px;
}


</style>

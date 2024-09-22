<script>

export default {
    
	props: ['params'],
    
	data: function() {
		return {
			datestr: this.$timestamp2date(this.params.Upload_time)
		}
	},
	methods: {
		remove() {
			this.$axios({
				method: 'delete',
				url: `/users/${this.$username.value}/comments/${this.$photo.value.PhotoId}/${this.params.Id}`,
				headers: {Authorization: this.$token.value},
			})
			.then(_ => {
				this.$emit('changed');
			})
			.catch(error => {
				console.error(error.response);
			});
		},
	},
}
</script>

<template>
    <div class="comment-container">
      <div class="comment-content">
        <p class="author">{{ params.Author }}</p>
        <p class="body">{{ params.Body }}</p>
      </div>
      <div class="comment-details">
        <p class="date">{{ this.datestr }}</p>
        <button v-if="params.Author == this.$username.value" type="button" class="btn btn-delete" @click="remove">
          Delete comment
        </button>
        <p>Deleting a comment for photo: {{ params.PhotoId }}</p>
      </div>
    </div>
  </template>
  

<style>
.comment-container {
  display: flex;
  justify-content: space-between;
  align-items: flex-start; /* Align items to the start */
  background-color: rgb(255, 243, 220);
  border-radius: 10px;
  margin: 0 auto;
  margin-top: 20px;
  padding: 20px;
  max-width: 800px;
}

.comment-content {
  flex: 1;
  margin-right: 20px; /* Space between the comment content and the buttons */
}

.author {
  margin-bottom: 5px;
  font-weight: bold;
}

.body {
  margin: 0;
}

.comment-details {
  text-align: right; /* Aligns the date and button to the right */
}

.date {
  margin-bottom: 5px;
  font-weight: bold;
}

.btn-delete {
  background-color: gray;
  color: white;
  border: 1px solid gray;
  border-radius: 4px;
  padding: 4px 8px; /* Narrow button */
  font-size: 12px; /* Smaller font size */
  font-weight: bold;
  text-transform: uppercase;
  cursor: pointer;
  transition: background-color 0.3s, border-color 0.3s;
}

.btn-delete:hover {
  background-color: darkgray;
  border-color: darkgray;
}
</style>

<script>
import PostComments from '../components/post_comments.vue'; // Adjust the path to your component
import DeleteComments from '../components/delete_comments.vue'; // Adjust the path to your component
import Likes from '../components/like_unlike.vue';  // Import the Likes component

export default {
    components: {
    PostComments,  // Register PostComments component
    DeleteComments, // Register DeleteComments component
    Likes  // Register the Likes component
  },
	props: ['params'],
	data: function() {
		return {
			datestr: this.$timestamp2date(this.params.Upload_time )
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
				url: `/users/${this.$username.value}/photos/${this.params.PhotoId }`,
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
    <div class="photo-container">
      <div class="photo-wrapper">
        <div v-for="photo in params.photodata" :key="photo.id">
        
        <img class="photo" :src="photo.png64" alt="User Photo" @click="view">
    
        <p class="caption">{{ photo.caption }}</p> <!-- Display the caption here -->
        <div class="info-container">
          <div class="info-left">
            <p><b>{{ datestr }}</b></p>
         
        
        <div class="info-right">
          <p><b>Comments:</b> {{ photo.Comment_count }}</p>
           <!-- Post Comments component -->
            <PostComments 
            :params="{
              PhotoId : photo.id,
              Photo_png:photo.png64,
              Caption:photo.caption,
              Upload_time:photo.timestamp,
              Comments:photo.comments,
              Comment_count:photo.comment_count,
              Like_count:photo.like_count,
              Likes:photo.likes
            }" />

            <!-- Delete Comments component -->
            <DeleteComments 
            :params="{
              PhotoId : photo.id,
              Photo_png:photo.png64,
              Caption:photo.caption,
              Upload_time:photo.timestamp,
              Comments:photo.comments,
              Comment_count:photo.comment_count,
              Like_count:photo.like_count,
              Likes:photo.likes
            }"/>

          <p><b>Likes:</b> {{ photo.Like_count  }}</p>
          <Likes :users="photo.Likes" />
          <button v-if="params.Author == $username.value" type="button" class="btn btn-delete" @click="remove">
            Delete
          </button>
        </div>
        </div>
    </div>
  </div>
   </div>
    </div>
  </template>
  

  <style>
  .photo-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    max-width: 1000px;
    background-color: rgb(212, 204, 188);
    border-radius: 10px;
    margin: 0 auto;
    margin-top: 20px;
    padding: 20px;
  }
  
  .photo-wrapper {
    margin-bottom: 10px; /* Space between photo and caption */
  }
  
  .photo {
    max-height: 50vh;
    max-width: 100%;
    height: auto;
    width: auto;
    border-radius: 5px;
    cursor: pointer; /* Pointer cursor on hover */
  }
  
  .caption {
    font-style: italic; /* Style for the caption text */
    margin-bottom: 20px; /* Space below the caption */
  }
  
  .info-container {
    display: flex;
    justify-content: space-between;
    width: 100%;
  }
  
  .info-left,
  .info-right {
    display: flex;
    flex-direction: column;
  }
  
  .btn-delete {
    background-color: gray;
    color: white;
    border: 1px solid gray;
    border-radius: 4px;
    padding: 5px 10px;
    font-size: 12px;
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
  

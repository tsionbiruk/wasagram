<script>
import FollowBan from '@/components/follow_ban.vue'; 
import Deletepost from '@/components/delete_post.vue'; 
export default {
    components: {
    FollowBan, 
    Deletepost,
  },
  data: function() {
    return {
      profil_pic: [],
      photodata: [],
      following: [],
      followers: [],
      banned: [],
      photo_count: [],
      errormsg: null,
      loading: false,
      form: {
        username: null,
      },
    };
  },
  methods: {

    searchUser() {
      this.refresh(); // Call refresh to load new profile data for the username in form.username
    },

    async refresh() {
      this.loading = true;
      this.errormsg = null;
      this.photodata = [];
      try {
        const response = await this.$axios({
          method: "get",
          url: `/users/${this.$username.value}/profile/${this.form.username}`,
          headers: { Authorization: this.$token.value },
        });

        const photos = response.data.Photo;
        this.profil_pic = response.data.ProfilPic;
        this.following = response.data.Following_count;
        this.followers = response.data.Follower_count;
        this.banned = response.data.banned;
        this.photo_count = response.data.Photo_count;

        for (let i = 0; i < photos.length; i++) {
          const { PhotoId: id, Caption: caption, timestamp: ts, comments, comment_count, likes, like_count } = photos[i];

          let png64 = "";
          let imageResponse = null;
          try {
            imageResponse = await this.$axios.get(`/photos/${id}`, {
              responseType: "arraybuffer",
            });
          } catch (error) {
            this.errormsg = `Failed to retrieve photo ${id}: ${error.response.data}`;
            console.error(error);
          }

          if (imageResponse && imageResponse.data) {
            let binary = '';
            const bytes = new Uint8Array(imageResponse.data);
            const len = bytes.byteLength;
            for (let i = 0; i < len; i++) {
                binary += String.fromCharCode(bytes[i]);
            }
            png64 = `data:image/png;base64,${btoa(binary)}`; // Prefix with data URL format
            }


          this.photodata.push({
            id: id,
            png64: png64,
            caption: caption,
            timestamp: ts,
            comments: comments,
            comment_count: comment_count,
            likes: likes,
            like_count: like_count,
          });
        }
      } catch (error) {
        this.errormsg = `Failed to retrieve profile data: ${error.response.data}`;
        console.error(error);
      } finally {
        this.loading = false;
      }
    },
    rename() {
      this.$axios({
        method: "post",
        url: `/users/${this.$username.value}/username`,
        data: JSON.stringify(this.form.username),
        headers: { Authorization: this.$token.value },
      })
        .then(() => {
          this.errormsg = null;
          this.$username.value = this.form.username;
          this.refresh();
        })
        .catch((error) => {
          this.errormsg = `User rename failed: ${error.response.data}`;
          console.error(error.response);
        });
    },
    upload() {
      let file = this.$refs.png.files[0];
      const reader = new FileReader();
      reader.onerror = (error) => {
        this.errormsg = `Failed to read photo file: ${error.message}`;
        console.error(error);
      };
      reader.onload = (res) => {
        this.$axios({
          method: "post",
          url: `/users/${this.$username.value}/photos`,
          data: res.target.result,
          headers: { Authorization: this.$token.value },
        })
          .then(() => {
            this.errormsg = null;
            this.refresh();
          })
          .catch((error) => {
            this.errormsg = `Photo upload failed: ${error.response.data}`;
            console.error(error.response);
          });
      };
      reader.readAsArrayBuffer(file);
    },
    triggerUpload() {
        // Trigger click on hidden file input
        this.$refs.png.click();
    }
  },
  mounted() {
    this.$photoorigin.value = "profile";
    this.refresh();
  },
};
</script>


<template>
    <div>
    <input
      type="text"
      v-model="form.username"   
      
      placeholder="Search for a user"
    />
    <button @click="searchUser">Search</button> <!-- Button to trigger search -->

    <div v-if="loading">Loading...</div>
    <div v-if="errormsg">{{ errormsg }}</div>

     <!-- Include FollowBan component here -->
     <FollowBan 
      :params="{
        username: form.username,
        ProfilPhoto: profil_pic // Assuming this is where the profile picture is stored
      }" 
      @changed="refresh" 
    />
    
    <div v-for="photo in photodata" :key="photo.id">
        
        <img class="photo" :src="photo.png64" alt="User Photo"/>
    
        <p class="caption">{{ photo.caption }}</p> 
   </div>
   <div>
        <div>
			<h5>My photo stream:</h5>
            <div class="d-flex flex-column flex-wrap flex-md-nowrap align-items-between justify-items-center" style="max-width: 1000px; background-color: rgb(212, 204, 188); border-radius: 10px; margin: 0 auto; margin-top: 20px; padding: 20px;">
                <h5>Upload photo:</h5>
                <!-- Hidden file input -->
                <input ref="png" @change="upload" type="file" style="display: none;" id="file-upload">
                
                <!-- Button to trigger file upload -->
                <button class="btn btn-primary" @click="triggerUpload">Choose File</button>
            </div>
        </div>
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
			

			
			<div v-if="this.photodata.length">
				<div v-for="photo in this.photodata" :key="photo.timestamp">
					<StreamPost :params="photo" @deleted="refresh"></StreamPost>
				</div>
			</div>
			<div v-else-if="this.$token.value">
				<p><br>Your own photo stream is empty!</p>
			</div>
		</div>
    </div>
    </div>
</template>

<style>
</style>
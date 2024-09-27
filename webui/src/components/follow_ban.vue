<script>
export default {
	props: ['params'],
	data: function() {
		return {
			isfollow: false,
			isban: false,
            isfollowing: false ,
            profilePhotoBase64: '',
		}
	},
	created() {
		this.refresh();
        this.convertProfilePhotoToBase64();
	},
	methods: {
        convertProfilePhotoToBase64() {
			if (this.params.ProfilPhoto) {
				// Assuming params.ProfilPhoto is in binary or can be converted to Base64
				const binary = String.fromCharCode(...new Uint8Array(this.params.ProfilPhoto));
				this.profilePhotoBase64 = btoa(binary); // Convert to Base64
			}
		},
		refresh() {
			if (!this.$token.value) {
				return;
			}
			this.$axios({
				method: 'get',
				url: `/users/${this.$username.value}/banned/${this.params.username}`,
				headers: {Authorization: this.$token.value},
			})
			.then(response => {
				let banned = response.data;
				for (let i = 0; i < banned.length; i++) {
					if (this.username == banned[i]) {
						this.isban = true;
						break;
					}
				}
			})
			.catch(error => {
				console.error(error.response);
			});

			this.$axios({
				method: 'get',
				url: `/users/${this.$username.value}/followers/${this.params.username}`,
				headers: {Authorization: this.$token.value},
			})
			.then(response => {
				let follower = response.data;
				for (let i = 0; i < following.length; i++) {
					if (this.username == follower[i]) {
						this.isfollow = true;
						break;
					}
				}
			})
			.catch(error => {
				console.error(error.response);
			});

            this.$axios({
				method: 'get',
				url: `/users/${this.$username.value}/following/${this.params.username}`,
				headers: {Authorization: this.$token.value},
			})
			.then(response => {
				let following = response.data;
				for (let i = 0; i < following.length; i++) {
					if (this.username == following[i]) {
						this.isfollowing = true;
						break;
					}
				}
			})
			.catch(error => {
				console.error(error.response);
			});
		},
		setfollow(state) {
			this.$axios({
				method: state ? 'put' : 'delete',
				url: `/users/${this.$username.value}/follow/${this.params.username}`,
				headers: {Authorization: this.$token.value},
			})
			.then(_ => {
				this.isfollow = state;
				this.$emit('changed');
			})
			.catch(error => {
				console.error(error.response);
			});
		},
		setban(state) {
			this.$axios({
				method: state ? 'put' : 'delete',
				url: `/users/${this.$username.value}/ban/${this.params.username}`,
				headers: {Authorization: this.$token.value},
			})
			.then(_ => {
				this.isban = state;
				this.$emit('changed');
			})
			.catch(error => {
				console.error(error.response);
			});
		},
        getfollower() {
            this.$axios({
				method: 'get',
				url: `/users/${this.$username.value}/followers/${this.params.username}`,
				headers: {Authorization: this.$token.value},
			})
            .then(response => {
            // Assuming the API returns an array of banned users in response.data
            return response.data;  // This will return the list of banned users
        })
            .catch(error => {
                console.error('Error fetching banned users:', error);
                throw error;  // You can choose how to handle the error
            });
        },

        getfollowing() {
            this.$axios({
				method: 'get',
				url: `/users/${this.$username.value}/following/${this.params.username}`,
				headers: {Authorization: this.$token.value},
			})
            .then(response => {
            // Assuming the API returns an array of banned users in response.data
            return response.data;  // This will return the list of banned users
        })
            .catch(error => {
                console.error('Error fetching banned users:', error);
                throw error;  // You can choose how to handle the error
            });
        },
        toggleFollow() {
        this.setfollow(!this.isfollow);  // Toggle the follow/unfollow state
        },

        toggleban() {
        this.setban(!this.isban);  
        },

    

        getbanned() {
            this.$axios({
				method: 'get',
				url: `/users/${this.$username.value}/banned/${this.params.username}`,
				headers: {Authorization: this.$token.value},
			})
            .then(response => {
            // Assuming the API returns an array of banned users in response.data
            return response.data;  // This will return the list of banned users
        })
            .catch(error => {
                console.error('Error fetching banned users:', error);
                throw error;  // You can choose how to handle the error
            });
        }

	}
}
</script>

<template>
<div id="users-info" class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center" style="max-width: 1000px; background-color: rgb(212, 204, 188) ; border-radius: 10px; margin: 0 auto; margin-top: 20px; padding: 20px;">
    <div class="user-info d-flex align-items-center">
        <img :src="profilePhotoBase64" alt="Profile Picture" style="width: 50px; height: 50px; border-radius: 50%; margin-right: 10px;">
 
             
        <span>{{ params.username }}</span>
    </div>
	<div v-if="this.$token.value && params.username != this.$username.value">
        <div class="button-group primary-group">
            <button type="button" 
                    class="btn btn-sm follow-btn" 
                    :class="isfollow ? 'unfollow-btn' : 'follow-btn'" 
                    @click="toggleFollow">
                {{ isfollow ? 'Unfollow' : 'Follow' }}
            </button>

            <button type="button" 
                    class="btn btn-sm follow-btn" 
                    :class="isban ? 'unfollow-btn' : 'follow-btn'" 
                    @click="toggleBan">
                {{ isban ? 'Unban' : 'Ban' }}
            </button>
        </div>

        <!-- Second button group: Followers, Following, Banned -->
        <div class="button-group secondary-group">
            <button type="button" class="btn btn-sm follow-btn" @click="getfollower()">
                Followers
            </button>
            <button type="button" class="btn btn-sm follow-btn" @click="getfollowing()">
                Following
            </button>
            <button type="button" class="btn btn-sm follow-btn" @click="getbanned()">
                Banned
            </button>
        </div>
	</div>
	<div v-else-if="this.$token.value && isban=== true">
		<p><b>Can not view profil. You are banned!</b></p>
    </div>
    
    <div v-else="this.$token.value && username === this.$username.value">
        <div class="button-group secondary-group">
            <button type="button" class="btn btn-sm follow-btn" @click="getfollower()">
                Followers
            </button>
            <button type="button" class="btn btn-sm follow-btn" @click="getfollowing()">
                Following
            </button>
            <button type="button" class="btn btn-sm follow-btn" @click="getbanned()">
                Banned
            </button>
        </div>
    </div>
</div>
</template>

<style>
#users-info {
    display: flex;
    flex-direction: column; /* Stack children vertically */
    align-items: center; /* Center align children */
    max-width: 1000px;
    background-color: rgb(212, 204, 188);
    border-radius: 10px;
    margin: 0 auto;
    margin-top: 20px;
    padding: 20px;
}

.user-info {
    display: flex;
    flex-direction: column; /* Stack image and username vertically */
    align-items: center; /* Center align image and username */
    margin-bottom: 10px; /* Space between profile picture and buttons */
}

/* Base button style for all buttons */
.follow-btn, .unfollow-btn, .ban-btn {
    border: none;
    width: 150px;  /* Same width for consistency */
    height: 40px;  /* Same height for consistency */
    border-radius: 20px;  /* Smooth rectangular shape */
    font-size: 14px;
    transition: background-color 0.3s ease;
    cursor: pointer;
}

/* Follow button styles */
.follow-btn {
    background-color: #097ba1;  /* Light blue */
    color: white;
}

.follow-btn:hover {
    background-color: #3058a1;  /* Darker blue on hover */
}

/* Unfollow button styles */
.unfollow-btn {
    background-color: #D3D3D3;  /* Light gray */
    color: white;
}

.unfollow-btn:hover {
    background-color: #A9A9A9;  /* Darker gray on hover */
}

/* First button group: Follow/Unfollow and Ban/Unban */
.primary-group {
    display: flex;
    justify-content: flex-start;
    gap: 40px; /* More space between the two buttons */
    margin-bottom: 20px; /* Space between first and second group */
}

/* Second button group: Followers, Following, Banned */
.secondary-group {
    display: flex;
    justify-content: flex-start;
    gap: 20px;  /* Less space between the three buttons */
}


</style>
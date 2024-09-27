<script>
export default {
    props: ['params'],

	data: function() {
		return {
			islike: false,
        }
	},
	created() {
		this.refresh();
	},

methods: {
    toggleLike() {
            if (this.islike) {
                this.unlike();
            } else {
                this.like();
            }
        },

    refresh() {
			for (let i = 0; i < params.Likes.length; i++) {
				if (this.$username.value == params.Likes[i]) {
					this.islike = true;
					break;
				}
			}

		},
    like() {
        if (this.islike) {
            
            alert('Photo already liked');
            return; 
        }
        this.$axios({
            method: 'put',
            url: `/users/${this.$username.value}/like/${params.PhotoId}`,
            headers: { Authorization: this.$token.value },
        })
        .then(response => {
            this.islike = true;
            this.$emit('changed');
        })
        .catch(error => {
            console.error(error.response);
        });
    },
    unlike() {

        if (!this.islike) {
            
            alert('Photo was never liked');
            return; 
        }
        this.$axios({
            method: 'delete',
            url: `/users/${this.$username.value}/unlike/${params.PhotoId}`,
            headers: { Authorization: this.$token.value },
        })
        .then(response => {
            this.islike = false;
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
    

        <!-- Container for the like icon and buttons -->
        <div class="like-container">
            <!-- Like icon -->
            <i :class="{'fa-heart': islike, 'fa-heart-o': !islike}" 
               class="fa" 
               @click="toggleLike"
               aria-hidden="true"></i>

            <!-- Like/Unlike Buttons -->
            <div class="button-group">
                <button v-if="this.islike" type="button" class="btn btn-unlike" @click="unlike">
                    Unlike
                </button>
                <button v-else type="button" class="btn btn-like" @click="like">
                    Like
                </button>

            </div>

            
            <div v-for="like in params.Likes" :key="like">
            <likes :params="like" @changed="this.$emit('changed')"></likes>
            </div>

            <div class="like-item">
            <p>{{ params }}</p>
            </div>

        </div>
    
</template>


<style>



/* Container for the like icon and buttons */
.like-container {
    margin-top: 20px; /* Adds space between the photo and the like section */
}
.like-item {
  margin: 5px;
  font-size: 14px;
  
}

.button-group {
    margin-top: 10px; /* Adds space between the icon and buttons */
}

.btn-like {
    background-color: red; /* Red background color */
    color: white; /* White text color */
    border: 1px solid red; /* Red border */
    border-radius: 4px; /* Rounded corners */
    padding: 8px 16px; /* Padding for rectangle shape */
    font-size: 14px; /* Font size */
    font-weight: bold; /* Bold text */
    text-transform: uppercase; /* Uppercase text */
    cursor: pointer; /* Pointer cursor on hover */
    transition: background-color 0.3s, border-color 0.3s; /* Smooth transition */
}

.btn-like:hover {
    background-color: darkred; /* Darker red on hover */
    border-color: darkred; /* Border color on hover */
}

.btn-unlike {
    background-color: gray; /* Gray background color */
    color: white; /* White text color */
    border: 1px solid gray; /* Gray border */
    border-radius: 4px; /* Rounded corners */
    padding: 8px 16px; /* Padding for rectangle shape */
    font-size: 14px; /* Font size */
    font-weight: bold; /* Bold text */
    text-transform: uppercase; /* Uppercase text */
    cursor: pointer; /* Pointer cursor on hover */
    transition: background-color 0.3s, border-color 0.3s; /* Smooth transition */
}

.btn-unlike:hover {
    background-color: darkgray; /* Darker gray on hover */
    border-color: darkgray; /* Border color on hover */
}
.fa-heart {
    color: red; /* Color for liked state */
    cursor: pointer;
}

.fa-heart-o {
    color: grey; /* Color for not liked state */
    cursor: pointer;
}
</style>
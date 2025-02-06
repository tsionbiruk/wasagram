<script>
export default {
  props: ['username'],
  data() {
    return {
      isFollow: false,
      isBan: false,
      errorMsg: '',
    };
  },
  created() {
    this.refresh();
  },
  methods: {
    async refresh() {
      if (!this.$token?.value) return;

      try {
        const bannedResponse = await this.$axios.get(`/users/${this.$username.value}/banned`, {
          headers: { Authorization: this.$token.value },
        });
        this.isBan = bannedResponse.data.includes(this.username);
      } catch (error) {
        console.error('Error fetching banned users:', error);
      }

      try {
        const followResponse = await this.$axios.get(`/users/${this.$username.value}/followed`, {
          headers: { Authorization: this.$token.value },
        });
        this.isFollow = followResponse.data.includes(this.username);
      } catch (error) {
        console.error('Error fetching followed users:', error);
      }
    },
    async setFollow(state) {
      try {
        await this.$axios({
          method: state ? 'put' : 'delete',
          url: `/users/${this.$username.value}/followed/${this.username}`,
          headers: { Authorization: this.$token.value },
        });
        this.isFollow = state;
        this.$emit('changed');
      } catch (error) {
        console.error('Error updating follow status:', error);
      }
    },
    async setBan(state) {
		try {
			await this.$axios({
			method: state ? 'put' : 'delete',
			url: `/users/${this.$username.value}/banned/${this.username}`,
			headers: { Authorization: this.$token.value },
			});
			this.isBan = state;
			this.$emit('changed');
		} catch (error) {
			console.error('Ban request failed:', error.response?.data || error.message);
			this.errorMsg = error.response
			? `Error ${error.response.status}: ${error.response.data?.message || error.response.statusText}`
			: 'An unexpected error occurred.';
		}
	},

  },
};
</script>

<template>
  <div id="users-entry" class="user-entry">
    <span>{{ username }}</span>
    <div v-if="$token?.value && username !== $username.value">
      <button v-if="isFollow" class="btn btn-sm btn-primary" @click="setFollow(false)">
        Unfollow
      </button>
      <button v-else class="btn btn-sm btn-outline-primary" @click="setFollow(true)">
        Follow
      </button>
      <button v-if="isBan" class="btn btn-sm btn-danger" @click="setBan(false)">
        Unban
      </button>
      <button v-else class="btn btn-sm btn-outline-danger" @click="setBan(true)">
        Ban
      </button>
    </div>
    <div v-else-if="$token?.value && username === $username.value">
      <p><b>You!</b></p>
    </div>
    <p v-if="errorMsg" class="error-message">{{ errorMsg }}</p>
  </div>
</template>

<style>
.user-entry {
  max-width: 1000px;
  background-color: rgb(212, 204, 188);
  border-radius: 10px;
  margin: 20px auto;
  padding: 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.error-message {
  color: red;
  margin-top: 10px;
}
</style>

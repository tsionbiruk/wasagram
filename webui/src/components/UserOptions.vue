<script>
export default {
	props: ['username'],
	data: function() {
		return {
			isfollow: false,
			isban: false,
		}
	},
	created() {
		this.refresh();
	},
	methods: {
		refresh() {
			if (!this.$token.value) {
				return;
			}
			this.$axios({
				method: 'get',
				url: `/users/${this.$username.value}/banned`,
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
				url: `/users/${this.$username.value}/followed`,
				headers: {Authorization: this.$token.value},
			})
			.then(response => {
				let following = response.data;
				for (let i = 0; i < following.length; i++) {
					if (this.username == following[i]) {
						this.isfollow = true;
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
				url: `/users/${this.$username.value}/followed/${this.username}`,
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
				url: `/users/${this.$username.value}/banned/${this.username}`,
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
	},
}
</script>

<template>
<div id="users-entry" class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center" style="max-width: 1000px; background-color: rgb(212, 204, 188) ; border-radius: 10px; margin: 0 auto; margin-top: 20px; padding: 20px;">
	<span>{{ username }}</span>
	<div v-if="this.$token.value && username != this.$username.value">
		<button v-if="this.isfollow" type="button" class="btn btn-sm btn-primary" @click="setfollow(false)">
			Unfollow
		</button>
		<button v-else type="button" class="btn btn-sm btn-outline-primary" @click="setfollow(true)">
			Follow
		</button>
		<button v-if="this.isban" type="button" class="btn btn-sm btn-danger" @click="setban(false)">
			Unban
		</button>
		<button v-else type="button" class="btn btn-sm btn-outline-danger" @click="setban(true)">
			Ban
		</button>
	</div>
	<div v-else-if="this.$token.value && username == this.$username.value">
		<p><b>You!</b></p>
	</div>
</div>
</template>

<style>
#users-entry {
	display: flex;
	flex-flow: row nowrap;
	justify-content: space-between;
}
</style>

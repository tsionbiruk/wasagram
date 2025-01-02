<script>
export default {
	data: function() {
		return {
			islike: false,
        }
	},
	created() {
		this.refresh();
	},
	methods: {
		refresh() {
			for (let i = 0; i < this.$photo.value.likes.length; i++) {
				if (this.$username.value == this.$photo.value.likes[i]) {
					this.islike = true;
					break;
				}
			}

		},
		setlike(state) {
			this.$axios({
				method: state ? 'put' : 'delete',
				url: `/users/${this.$username.value}/likes/${this.$photo.value.id}`,
				headers: {Authorization: this.$token.value},
			})
			.then(_ => {
				this.islike = state;
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
	<button v-if="this.islike" type="button" class="btn btn-sm btn-primary" @click="setlike(false)">
		Unlike
	</button>
	<button v-else type="button" class="btn btn-sm btn-outline-primary" @click="setlike(true)">
		Like
	</button>
</template>

<style>
</style>

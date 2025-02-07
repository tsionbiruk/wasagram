<script>
export default {
	data: function() {
		return {
			users: [],
			errormsg: null,
			loading: false,
		}
	},
	methods: {
		async refresh() {
			this.loading = true;
			this.errormsg = null;
			this.$axios({
				method: 'get',
				url: '/users',
			})
			.then(async response => {
				this.errormsg = null;
				this.users = response.data;
			})
			.catch(error => {
				this.errormsg = `Failed to retrieve user list: ${error.response.data}`
				console.error(error);
			});
			this.loading = false;
		},
	},
	mounted() {
		this.refresh();
	}
}
</script>

<template>
	<div>
		<div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
			<h1 class="h2">All users {{ this.users.length }}</h1>
			<div v-if="this.$token.value" class="btn-toolbar mb-2 mb-md-0">
				<div class="btn-group me-2">
					<button type="button" class="btn btn-sm btn-outline-secondary" @click="refresh">
						Refresh
					</button>
				</div>
			</div>
		</div>

		<ol>
			<div v-for="name in this.users" :key="name">
				<UserOptions :username="name"></UserOptions>
			</div>
		</ol>

		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
	</div>
</template>

<style>
</style>

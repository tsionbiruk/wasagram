<script>
export default {
	data() {
		return {
			users: [], 
			searchQuery: '', 
			errormsg: null,
			loading: false,
		};
	},

	computed: {
		filteredUsers() {
			if (!this.searchQuery.trim()) {
				return this.users; 
			}
			return this.users.filter(user =>
				user.toLowerCase().includes(this.searchQuery.toLowerCase())
			);
		},
	},

	methods: {
		async refresh() {
			this.loading = true;
			this.errormsg = null;
			try {
				const response = await this.$axios.get('/users', {
					headers: { Authorization: this.$token.value }
				});
				this.users = response.data; // Store retrieved users
			} catch (error) {
				this.errormsg = `Failed to retrieve user list: ${error.response?.data || 'Unknown error'}`;
				console.error(error);
			}
			this.loading = false;
		},
	},

	mounted() {
		this.refresh(); // Load users on page mount
	},
};
</script>


<template>
	<div>
		<!-- 🔎 Search Bar -->
		<div class="search-container">
			<input 
				v-model="searchQuery"
				class="search-box"
				type="text"
				placeholder="Search username..."
			/>
		</div>

		<!-- 🔹 Page Header -->
		<div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
			<h1 class="h2">All users ({{ filteredUsers.length }})</h1>
			<div v-if="$token.value" class="btn-toolbar mb-2 mb-md-0">
				<div class="btn-group me-2">
					<button type="button" class="btn btn-sm btn-outline-secondary" @click="refresh">
						Refresh
					</button>
				</div>
			</div>
		</div>

		<!-- 🔹 User List -->
		<ol>
			<div v-for="name in filteredUsers" :key="name">
				<UserOptions :username="name"></UserOptions>
			</div>
		</ol>

		<!-- 🔹 Error Message -->
		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
	</div>
</template>


<style scoped>
.search-container {
	margin-bottom: 0px;
	text-align: center;
}

.search-box {
	width: 60%;
	max-width: 400px;
	padding: 8px;
	border: 1px solid #ccc;
	border-radius: 5px;
	font-size: 16px;
	margin-top: 20px;
}

.user-list {
	display: flex;
	flex-wrap: wrap;
	gap: 10px;
}
</style>


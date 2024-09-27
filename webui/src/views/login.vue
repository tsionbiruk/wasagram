<script>
export default {
	data: function() {
		return {
			errormsg: null,
			form: {
				username: '',
			},
		}
	},
	methods: {
		login(event) {
			event.preventDefault(); // Prevent form submission refresh
			this.$axios({
				method: 'post',
				url: `/session`,
				data: this.form, // Pass the form data as an object
			})
			.then(response => {
				this.errormsg = null;
				this.$token.value = response.data;
				this.$username.value = this.form.username;
				this.$router.push(`/profile/${this.$username.value}`);

				console.log(`successfully logged in with token ${this.$token.value}`);
			})
			.catch(error => {
				this.errormsg = `Login failed: ${error.response.data}`;
				console.error(error.response);
			});
		},
	},
}
</script>

<template>
	<div>
		<div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
			<h1 class="h2">Log-in</h1>
		</div>
		<div v-if="!this.$token.value">
			<form v-on:submit.prevent="login"> <!-- Added prevent default -->
				<label for="form-username">Username</label>
				<br>
				<input v-model="form.username" id="form-username">
				<button class="btn btn-sm btn-outline-secondary" type="submit">Log In</button>
			</form>
		</div>
		<div v-else>
			<p>Logged in as {{ this.$username.value }}.</p>
		</div>

		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
	</div>
</template>

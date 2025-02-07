<script>
export default {
	props: ['params'],
	data: function() {
		return {
			datestr: this.$timestamp2date(this.params.timestamp)
		}
	},
	methods: {
		remove() {
			this.$axios({
				method: 'delete',
				url: `/users/${this.$username.value}/comments/${this.$photo.value.id}/${this.params.id}`,
				headers: {Authorization: this.$token.value},
			})
			.then(_ => {
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
<div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center" style="max-width: 800px; background-color: rgb(141, 164, 207) ; border-radius: 10px; margin: 0 auto; margin-top: 20px; padding: 20px;">
	<div>
		<p style= "margin-bottom: 5px; font-weight: bold;" >{{ params.user }}</p>
		<p>{{ params.content }}</p>
	</div>
	<div>
		<p style= "margin-bottom: 5px; font-weight: bold;" >{{ this.datestr }}</p>
		<button v-if="params.user == this.$username.value" size="sm" type="button" class="btn btn-sm btn-outline-danger" @click="remove">
			Delete comment
		</button>
	</div>
</div>
</template>

<style>
</style>

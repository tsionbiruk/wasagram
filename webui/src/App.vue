<script setup>
</script>
<script>
export default {
	created() {
		if (!this.$token.value) {
			this.$router.push('/login');
		}
	},
	methods: {
		logout() {
			this.$token.value = 0;
			this.$router.push('/login');
		},
	},
}
</script>

<template>

	<header class="navbar navbar-dark sticky-top bg-dark-blue flex-md-nowrap p-0 shadow">
		<a class="navbar-brand wasagram-header" href="#/">Wasagram</a>
			<button class="navbar-toggler position-absolute d-md-none collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#sidebarMenu" aria-controls="sidebarMenu" aria-expanded="false" aria-label="Toggle navigation">
				<span class="navbar-toggler-icon"></span>
			</button>
	</header>


	<div class="container-fluid">
		<div class="row">
			<nav v-if="this.$token.value" id="sidebarMenu" class="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse">
				<div class="position-sticky pt-3 sidebar-sticky">
					<ul class="nav flex-column">
						<li class="nav-item">
							<RouterLink to="/" class="nav-link sidebar-link">
								Photo Stream
							</RouterLink>
						</li>
						<li class="nav-item">
							<RouterLink to="/profile" class="nav-link sidebar-link">
								My Profile
							</RouterLink>
						</li>
						<li class="nav-item">
							<RouterLink to="/users" class="nav-link sidebar-link">
								All Users
							</RouterLink>
						</li>
						<li>
							<button type="button" style="margin: 10px 0px 0px 15px" class="btn btn-sm btn-outline-secondary" @click="logout">
								{{ this.$username.value }}: Log out
							</button>
						</li>
					</ul>
				</div>

			</nav>
			<nav v-else id="sidebarMenu" class="col-md-3 col-lg-2 d-md-block bg-beige sidebar collapse">
				<div class="position-sticky pt-3 sidebar-sticky">
					<h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
						<span>Not Logged In</span>
					</h6>
					<ul class="nav flex-column">
						<li class="nav-item">
							<RouterLink to="/login" class="nav-link">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#key"/></svg>
								Log In
							</RouterLink>
						</li>
						<li class="nav-item">
							<RouterLink to="/users" class="nav-link">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#layout"/></svg>
								Users
							</RouterLink>
						</li>
					</ul>
				</div>
			</nav>

			
			<main class="col-md-9 ms-sm-auto col-lg-10 px-md-4">
				<RouterView />
			</main>
		</div>
	</div>
</template>

<style>
#sidebarMenu {
  background-color: rgba(141, 156, 186, 0.894) !important;
}

.bg-dark-blue {
  background-color: #001f3f !important; 
}


.sidebar-link {
  font-size: 15px; 
  font-weight: bold; 
}


.sidebar-link:hover {
  font-size: 20px;
  transition: font-size 0.2s ease-in-out;
}


.wasagram-header {
  font-size: 24px; 
  
  padding: 10px 20px; 
}

</style>

<script>
	import Router,{replace} from "svelte-spa-router";

	import Main from "./views/Main.svelte";
	import Edit from "./views/Edit.svelte";

	const routes = {
		'/':Edit,
		'/edit':Edit,
	}

</script>

<main class="container-fluid">
	<Router {routes}></Router>
</main>

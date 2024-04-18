<svelte:options customElement="index-todo" />

<script context="module">
	var isLogout = false;
    export function redirectTo(url) {
      window.location.href = url;
    }
	export async function isAuthenticated() {
		const isAuth = document.cookie.includes("cookie");

		if (isAuth && !isLogout) {
			const cookieValue = getCookieValue("cookie");
			if (cookieValue != "") {
				var result = confirm("Vous êtes toujours connecté ! Voulez vous vous déconnecter ?")
				if (result) {
					isLogout = true;
					await logout();
				}
				else {
					redirectTo('app.html');
				}
			}
		}
	}

	function getCookieValue(cookieName) {
    const cookies = document.cookie.split("; ");
    for (const cookie of cookies) {
        const [name, value] = cookie.split("=");
        if (name === cookieName) {
            return value;
        }
    }
    return null;
}
</script>

<script>
	import { afterUpdate } from "svelte";
    import { logout } from "./app.svelte";

	afterUpdate (() => {
		isAuthenticated();
	})
</script>

<div id="index" class="centered">
    <h1>My TodoList</h1>
    <button on:click={() => redirectTo("signin.html")}>Inscription</button>
    <button on:click={() => redirectTo("login.html")}>Connexion</button>
</div>

<style>
	@import './style/style.css';	
	#index button {
		font-size: large;
	}
</style>
<svelte:options customElement="index-todo" />

<script context="module">
	let isLogout = false;
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
					await logout();
					isLogout = true
				}
				else {
					redirectTo('app.html');
				}
			}
			console.log(isLogout)
			return isLogout;
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
	import { onMount } from "svelte";
    import { logout } from "./app.svelte";

	onMount (() => {
		isAuthenticated();
	})
</script>

<style>
    h1{
		font-size: 70px;
	}
	
    .centered {
		width: 30em;
		margin: auto;
		display: inline-flexbox;
	}
	button {
    margin: 2%;
		height: 35px;
		width: 350px;
		border: none;
		border-radius: 5%;
		background-color: white;
		margin-left: 3%;
    font-size: large;
		transition: transform 0.2s ease;
	}
	button:hover{
		background-color: rgba(146, 146, 146, 0.181);
		cursor: pointer;
	}
  button:active {
		transform: scale(0.95);
	}
</style>
  
<div class="centered">
    <h1>My TodoList</h1>
    <button on:click={() => redirectTo("signin.html")}>Inscription</button>
    <button on:click={() => redirectTo("login.html")}>Connexion</button>
</div>
  
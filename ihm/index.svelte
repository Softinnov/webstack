<svelte:options customElement="index-todo" />

<script context="module">
    export function redirectTo(url) {
      window.location.href = url;
    }
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
</style>
  
<div class="centered">
    <h1>My TodoList</h1>
    <button on:click={() => redirectTo("signin.html")}>Inscription</button>
    <button on:click={() => redirectTo("login.html")}>Connexion</button>
</div>
  
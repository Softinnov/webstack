<svelte:options customElement="index-todo" />

<script>
    function redirectToSignIn() {
      window.location.href = 'signin.html';
    }
  
    function redirectToLogin() {
      window.location.href = 'login.html';
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
    <button on:click={redirectToSignIn}>Inscription</button>
    <button on:click={redirectToLogin}>Connexion</button>
</div>
  
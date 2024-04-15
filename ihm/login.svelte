<svelte:options customElement="login-todo" />

<script context="module">
    import { redirectTo, isAuthenticated } from './index.svelte';
    import { sendUser } from './signin.svelte';
    import { beforeUpdate } from "svelte";

    let email = '';
    let password = '';
  
    const handleSubmit = async () => {
        let user = {
            email,
            password
        }
        try{
            await sendUser(user,"login");
        }catch (error) {
            console.error(`Erreur lors de la connection au serveur : ${error.message}`);
        }
    };
</script>

<script>
	beforeUpdate (() => {
		isAuthenticated();
	})
</script>
  
<div class="centered">
    <h2>My TodoList</h2>
    <h3>Connexion :</h3>
    <form on:submit|preventDefault={handleSubmit}>
        <label>
        Email:
            <input type="email" style="margin-left: 106.22px;" bind:value={email} required>
        </label>
        <br>
        <label>
        Mot de passe:
            <input type="password" style="margin-left: 60px;" bind:value={password} required>
        </label>
  
        <button type="submit">Se connecter</button>
    </form>
    <div class="module">
        <p>Pas encore inscrit ?</p>
        <button on:click={() => redirectTo("signin.html")}>S'inscrire</button>
    </div>
</div>

<style>
    @import './style/style.css';
</style>

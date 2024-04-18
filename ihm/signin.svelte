<svelte:options customElement="signin-todo" />

<script context="module">
    import { customQueryEscape } from './app.svelte';
    import { redirectTo,isAuthenticated } from './index.svelte';
    import { answerResponse } from './app.svelte';

    let email = '';
    let password = '';
    let confirmpassword = '';
  
    const handleSubmit = async () => {
        let newuser = {
            email,
            password,
            confirmpassword
        }
        try {
            await sendUser(newuser,"signin");
        } catch (error) {
			console.error(`Erreur lors de la connection au serveur : ${error.message}`);
		}
    };

    export async function sendUser(user, route) {
		user.email = customQueryEscape(user.email);
        user.password = customQueryEscape(user.password);
        if (route == "signin"){
            user.confirmpassword = customQueryEscape(user.confirmpassword);
            var url =`/${route}?email=${user.email}&password=${user.password}&confirmpassword=${user.confirmpassword}`;
        } else {
            url = `/${route}?email=${user.email}&password=${user.password}`;
        }
        
		try {
			const reponse = await fetch(url,{method: "POST"});
			if (!reponse.ok) {
				const errorData = await reponse.text();
				answerResponse(errorData,reponse.status);
				throw new Error(`Erreur lors de la requête : ${reponse.status} ${reponse.statusText}`);
			}
			await reponse.text();
            if (reponse.ok){      
                redirectTo('app.html');
            }
		} catch (error) {
			console.error(error.message);
		}
	}
</script>
    
<script>
    import { beforeUpdate } from "svelte";

    beforeUpdate (() => {
        isAuthenticated();
    })
</script>

<div class="centered">

    <h2>My TodoList</h2>

    <h3>Inscription :</h3>
    
    <form on:submit|preventDefault={handleSubmit}>
        <label>
        Email:
            <input style="margin-left: 140.64px;" type="email" bind:value={email} required>
        </label>
        <br>
        <label>
        Mot de passe:
            <input style="margin-left: 94.42px;" type="password" bind:value={password} required>
        </label>
        <br>
        <label>
        Confirmer le mot de passe:
            <input style="margin-left: 10px;" type="password" bind:value={confirmpassword} required>
        </label>
        <br>
        <button type="submit">S'inscrire</button>
    </form>
    <p>Déjà inscrit ?</p>
    <button on:click={() => redirectTo("login.html")}>Se connecter</button>
</div>

<style>
    @import './style/style.css';
</style>
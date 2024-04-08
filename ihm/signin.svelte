<svelte:options customElement="signin-todo" />

<script context="module">
    import { customQueryEscape } from './app.svelte';
    import { redirectTo } from './index.svelte';

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
				alert(errorData);
				throw new Error(`Erreur lors de la requête : ${reponse.status} ${reponse.statusText}`);
			}
			await reponse.text();
            if (reponse.ok){      
                redirectTo('app.html');
            }
		} catch (error) {
			console.error(`Une erreur s'est produite : ${error.message}`);
		}
	}
</script>

<div class="centered">

    <h2>My TodoList</h2>

    <h3>Inscription :</h3>
    
    <form on:submit|preventDefault={handleSubmit}>
        <label>
        Email:
            <input type="email" bind:value={email} required>
        </label>
        <br>
        <label>
        Mot de passe:
            <input type="password" bind:value={password} required>
        </label>
        <br>
        <label>
        Confirmer le mot de passe:
            <input type="password" bind:value={confirmpassword} required>
        </label>
        <br>
        <button type="submit">S'inscrire</button>
    </form>
    <p>Déjà inscrit ?</p>
    <button on:click={() => redirectTo("login.html")}>Se connecter</button>
</div>

<style>
    h2{
		font-size: 50px;
	}

	label{
		font-size: medium;
	}
    .centered {
		width: 30em;
		margin: auto;
		display: grid;
	}
	p{
		font-size: small;
	}
    input{
        margin: 1%;
    }

	button {
        margin: 2%;
		height: 35px;
		width: 350px;
		border: none;
		border-radius: 5%;
		background-color: white;
		margin-left: 3%;
	}
	button:hover{
		background-color: rgba(146, 146, 146, 0.381);
		cursor: pointer;
	}
</style>
<svelte:options customElement="login-todo" />

<script>
    import { redirectTo } from './index.svelte';
    import { sendUser } from './signin.svelte'

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

  
<div class="centered">
    <h2>My TodoList</h2>
    <h3>Connexion :</h3>
    <form on:submit|preventDefault={handleSubmit}>
        <label>
        Email:
            <input type="email" style="margin-left: 66.3px;" bind:value={email} required>
        </label>
        <br>
        <label>
        Mot de passe:
            <input type="password" style="margin-left: 20px;" bind:value={password} required>
        </label>
  
        <button type="submit">Se connecter</button>
    </form>
    <div class="module">
        <p>Pas encore inscrit ?</p>
        <button on:click={() => redirectTo("signin.html")}>S'inscrire</button>
    </div>
</div>

<style>
    form{
        margin: 2%;
    }
    h2{
		font-size: 50px;
	}
    input{
        margin: 1%;
    }
	label{
        margin-left: 2%;
		font-size: medium;
	}
	p{
		font-size: small;
	}

	button {
        margin: 2%;
		height: 35px;
		width: 350px;
		border: none;
		border-radius: 5%;
		background-color: white;
		margin-left: 3%;
		transition: transform 0.2s ease;
	}
	button:hover{
		background-color: rgba(146, 146, 146, 0.381);
		cursor: pointer;
	}
    button:active {
		transform: scale(0.95);
	}
    .module{
        align-items: center;
        align-content: center;
    }
	.centered {
		width: 30em;
		margin: auto;
		display:grid;
	}
</style>

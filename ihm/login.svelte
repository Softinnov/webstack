<svelte:options customElement="login-todo" />

<script>
    import { sendUser } from './signin.svelte'
    // import { redirectToTodo } from './index.svelte'

    let email = '';
    let password = '';
  
    const handleSubmit = () => {
        let user = {
            email,
            password
        }
        try{
            sendUser(user,"login")
        }catch (error) {
            console.error(`Erreur lors de la connection au serveur : ${error.message}`);
        }
    };
</script>

  
<div>
    <h2>Connexion</h2>
    <form on:submit|preventDefault={handleSubmit}>
        <label>
        Email:
            <input type="email" bind:value={email} required>
        </label>
  
        <label>
        Mot de passe:
            <input type="password" bind:value={password} required>
        </label>
  
        <button type="submit">Se connecter</button>
    </form>
</div>
  
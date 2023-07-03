<template lang="pug">
b-container
 
  b-row
    b-col.mt-2
      b-img(:src="image" fluid )
  b-row.mt-3
    b-col
      b-alert(show variant="success")
        h1 READY TO CREATE YOUR OWN CHALLENGE? 
  b-row.mt-3
    b-col
      b-button(variant="success" pill size="lg" @click="connect")
        b-img(:src="metamask" rounded="circle" fluid height="60em" width="60em")
        | CONNECT WALLET
</template>

<script>
// @ is an alias to /src
import image from "../assets/onboarding.png";
import metamask from "../assets/metamask-icon.png";


export default {
  name: 'Home',
  components: {
  },
  data(){
    return {
        image,
        metamask,
        account:""
    }
  },
  methods:{
    connect(){
        if (typeof window.ethereum !== "undefined") {
        ethereum
         .request({ method: "eth_requestAccounts" })
         .then((accounts) => {
           const account = accounts[0]
           this.account = account
           this.$router.push({ name: 'create', params: {account: account }})


        }).catch((error) => {
           console.log(error, error.code);

        
       });
   } else {
       window.open("https://metamask.io/download/", "_blank");
   }
    }
  }
}
</script>

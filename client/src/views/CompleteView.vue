<template lang="pug">
b-overlay(:show="loading" rounded="sm")
     b-container
      b-row
        b-col
          b-navbar( type="dark" variant="info")
            b-navbar-brand(href="#") Comolete Challenge 
            b-navbar-nav
              b-nav-item-dropdown(:text="account_text")
                b-dropdown-item  My Challenges
                b-dropdown-item Complete Challenge
      b-row.mt-2
        b-alert(show variant="success")
          h4.alert-heading Challenge description 
          p  {{  description}}  
      b-row.mt-2
        b-col
          h3 Input photo to complete challenge
      b-row.mt-2 

        b-col
            b-card(:img-src="image" img-top )
              b-card-text
                b-form-file(v-model="file"  :file-name-formatter="formatNames" plain)
                b-button(variant="danger" @click="cleanFile")
                  b-icon(icon="x-octagon")
    
      b-row.mt-2 
          b-button(block variant="success" size="lg" @click="recognize")
            b-icon(icon="arrow-clockwise")
            | Recognize
          
      b-row.mt-2
        b-alert(:show="result==='ok'"  variant="success")
          h4.alert-heading Well done!
          hr
          h4  Congrats! You completed challenge
              
        b-alert(:show="result==='bad'"  variant="warning")
          h4.alert-heading Opps!
          hr
          h4 Try another photo
         
    
    
    
</template>
    
    <script>
    import {v4 as uuidv4} from 'uuid';
    
    import { BrowserProvider } from "ethers";
    import { WebKwil, Utils } from 'kwil';
    export const kwil = new WebKwil({
      kwilProvider: "https://provider.kwil.com",
    });
    const deployerAddress = "0x512F74323f28f2c0451DAe874E3EAfAC43Da2B69"
    const dbName = "photo_challenge"
    const dbid = Utils.generateDBID(deployerAddress, dbName);
    const provider = new BrowserProvider(window.ethereum);
    const signer = await provider.getSigner();
    
    
    
    const toBase64 = file => new Promise((resolve, reject) => {
        const reader = new FileReader();
        reader.readAsDataURL(file);
        reader.onload = () => resolve(reader.result);
        reader.onerror = reject;
    });
    
    
   
    export default {
      name: 'Complete',
      props:["account","challenge"],
      data(){
        return {
          file: null,
          fileB64: "",
          faces:[],
          loading:false,
          description:"some challenge",
          result: ""
        
        }
      },
      methods:{
        formatNames(){
          return ""
        },
        cleanFile(){
          this.file=null,
          this.fileB64 = ""
        },
      
        async recognize(){
          this.loading = true
          const resp = await fetch("/api/v1/recognize",{
            method:"POST",
            body:JSON.stringify({
              faces:[this.file1B64, this.file2B64, this.file3B64]
            }),
          })
        const data = await resp.json()
        if (resp.status !==200) {
            console.log("Bad status code")
        } else {
          this.faces = data.descriptors
    
        }
        this.recognizeButtonText = `Recognized ${Object.keys(this.faces).length} faces`
        this.loading = false
        this.recognzed = true
        },
        async createChallenge(){
          this.loading = true
          const inputs = new Utils.ActionInput()
          let recordId = uuidv4();
          inputs.put("$id",recordId)
          inputs.put("$faces_data", JSON.stringify(this.faces))
          inputs.put("$desription", this.desription)
          inputs.put("$completed",0)
          inputs.put("$complete_account","")
    
          const actionTx = await kwil
            .actionBuilder()
            .dbid(dbid)
            .name("create_challenge")
            .concat(inputs)
            .signer(signer)
            .buildTx();
           const tx= await kwil.broadcast(actionTx)
          
          this.loading = false
          this.hash = tx.data.txHash
          this.recordId = recordId
          this.saved = true
        }
      },
    
      computed:{
        account_text(){
          if (this.account){
            return `${this.account.substring(0,6)}...${this.account.substring(38)}`
          } else {
            return "No account"
          }
        },
         image(){
          if (this.fileB64 !=="" ){
            return this.fileB64
          } else {
            return ""
    
          }
        },
      
    
    
      },
      watch:{
        file : async function(){
          this.fileB64 = await toBase64(this.file)
    
        },
      },
      async mounted(){
        if (this.account === "" || this.account ===undefined) {
            const url = new URL(window.location.href)
            let path = url.pathname.split("/");
            
            this.$router.push({ name: 'home', params: {challenge: path.pop() }})

        } else {
            const inputs = new Utils.ActionInput()
            inputs.put("$id",this.challenge)

            const actionTx = await kwil
              .actionBuilder()
              .dbid(dbid)
              .name("get_challenge")
              .concat(inputs)
              .signer(signer)
              .buildTx();
            const tx = await kwil.broadcast(actionTx);
            if (tx.status !==200){
                print("Bad status: " + tx.status)
            } else {
                this.faces = tx.data.body[0].faces_data
                this.description = tx.data.body[0].description
            }
            

        }
        


      }
    
    }
    </script>
    
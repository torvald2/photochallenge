<template lang="pug">
b-overlay(:show="loading" rounded="sm")
 b-container
  b-row
    b-col
      b-navbar( type="dark" variant="info")
        b-navbar-brand(href="#") Create Challenge 
        b-navbar-nav
          b-nav-item-dropdown(:text="account_text")
            b-dropdown-item  My Challenges
            b-dropdown-item Complete Challenge
  b-row.mt-2
    b-col(cols="5")
      h1 Select up to 3 images
  b-row.mt-2
    b-col(cols="4")
        b-card(:img-src="image1" img-top )
          b-card-text
            b-form-file(v-model="file1"  :file-name-formatter="formatNames" plain)
            b-button(variant="danger" @click="cleanFile1")
              b-icon(icon="x-octagon")

    b-col(cols="4")
        b-card(:img-src="image2" img-top )
          b-card-text 
            b-form-file(  v-model="file2" :file-name-formatter="formatNames" plain)
            b-button(variant="danger" @click="cleanFile2")
              b-icon(icon="x-octagon")
    b-col(cols="4")
        b-card(:img-src="image3" img-top )
          b-card-text 
            b-form-file( v-model="file3"  :file-name-formatter="formatNames" plain)
            b-button(variant="danger" @click="cleanFile3")
              b-icon(icon="x-octagon")

  b-row.mt-2 
      b-button(block variant="success" size="lg" @click="recognize")
        b-icon(icon="arrow-clockwise")
        |  {{ recognizeButtonText }}
      
  b-row.mt-2(v-if="recognzed")
    b-col
      b-form-textarea(size="lg" v-model="description" placeholder="Put description. For example 'Find the photo where Elzabeth ii and Biden are together'")

  b-row.mt-2(v-if="recognzed")
      b-button(block variant="primary" size="lg" @click="createChallenge")
        |  CREATE CHALLENGE
  b-row.mt-2
    b-alert(:show="saved"  variant="success")
      h4.alert-heading Well done!
      p  Aww yeah, you successfully created challenge. TX hash {{ hash }}
      hr
      p 
        | Let`s chare this challenge to your friends 
        a(:href="complete_link") {{ complete_link }}



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


// @ is an alias to /src
let defaultImg = 'data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAn4AAAHECAYAAACnaSfuAAAAGXRFWHRTb2Z0d2FyZQBBZG9iZSBJbWFnZVJlYWR5ccllPAAAAyRpVFh0WE1MOmNvbS5hZG9iZS54bXAAAAAAADw/eHBhY2tldCBiZWdpbj0i77u/IiBpZD0iVzVNME1wQ2VoaUh6cmVTek5UY3prYzlkIj8+IDx4OnhtcG1ldGEgeG1sbnM6eD0iYWRvYmU6bnM6bWV0YS8iIHg6eG1wdGs9IkFkb2JlIFhNUCBDb3JlIDUuMy1jMDExIDY2LjE0NTY2MSwgMjAxMi8wMi8wNi0xNDo1NjoyNyAgICAgICAgIj4gPHJkZjpSREYgeG1sbnM6cmRmPSJodHRwOi8vd3d3LnczLm9yZy8xOTk5LzAyLzIyLXJkZi1zeW50YXgtbnMjIj4gPHJkZjpEZXNjcmlwdGlvbiByZGY6YWJvdXQ9IiIgeG1sbnM6eG1wPSJodHRwOi8vbnMuYWRvYmUuY29tL3hhcC8xLjAvIiB4bWxuczp4bXBNTT0iaHR0cDovL25zLmFkb2JlLmNvbS94YXAvMS4wL21tLyIgeG1sbnM6c3RSZWY9Imh0dHA6Ly9ucy5hZG9iZS5jb20veGFwLzEuMC9zVHlwZS9SZXNvdXJjZVJlZiMiIHhtcDpDcmVhdG9yVG9vbD0iQWRvYmUgUGhvdG9zaG9wIENTNiAoTWFjaW50b3NoKSIgeG1wTU06SW5zdGFuY2VJRD0ieG1wLmlpZDo3Q0QwMTIxQzlGN0UxMUUyOTJGM0JBRkQyNTI4N0VDNCIgeG1wTU06RG9jdW1lbnRJRD0ieG1wLmRpZDo3Q0QwMTIxRDlGN0UxMUUyOTJGM0JBRkQyNTI4N0VDNCI+IDx4bXBNTTpEZXJpdmVkRnJvbSBzdFJlZjppbnN0YW5jZUlEPSJ4bXAuaWlkOjI1OTdDMjYxOUY3QzExRTI5MkYzQkFGRDI1Mjg3RUM0IiBzdFJlZjpkb2N1bWVudElEPSJ4bXAuZGlkOjI1OTdDMjYyOUY3QzExRTI5MkYzQkFGRDI1Mjg3RUM0Ii8+IDwvcmRmOkRlc2NyaXB0aW9uPiA8L3JkZjpSREY+IDwveDp4bXBtZXRhPiA8P3hwYWNrZXQgZW5kPSJyIj8+LRfTYAAAK/hJREFUeNrs3X9s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'
export default {
  name: 'Create',
  props:['account'],
  data(){
    return {
      file1: null,
      file2: null, 
      file3: null,
      file1B64: "",
      file2B64: "",
      file3B64: "",
      faces:[],
      loading:false,
      recognizeButtonText: "Recognize faces",
      recognzed: false,
      description:"",
      hash:"",
      recordId:"",
      saved:false,
    }
  },
  methods:{
    formatNames(){
      return ""
    },
    cleanFile1(){
      this.file1=null,
      this.file1B64 = ""
    },
    cleanFile2(){
      this.file2=null,
      this.file2B64 = ""
    },
    cleanFile3(){
      this.file3=null,
      this.file3B64 = ""
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
      inputs.put("$desription", this.description)
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
      return `${this.account.substring(0,6)}...${this.account.substring(38)}`
    },
     image1(){
      if (this.file1B64 !=="" ){
        return this.file1B64
      } else {
        return defaultImg

      }
    },
    image2(){
      if (this.file2B64 !=="" ){
        return this.file2B64
      } else {
        return defaultImg

      }
    },
    image3(){
      if (this.file3B64 !=="" ){
        return this.file3B64
      } else {
        return defaultImg

      }
    },
    complete_link(){
      return `https://localhost:8080/complete/${this.recordId}`
    }


  },
  watch:{
    file1 : async function(){
      this.file1B64 = await toBase64(this.file1)

    },
    file2 : async function(){
      this.file2B64 = await toBase64(this.file2)

    },
    file3 : async function(){
      this.file3B64 = await toBase64(this.file3)

    }
  }

}
</script>

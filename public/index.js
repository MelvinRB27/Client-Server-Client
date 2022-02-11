let socket

new Vue({
    el: '#chat-app',
    data: {
        username:'',
        online: 0,
        target: '',
        file: '',
        disabled: false,
    },
    methods: {     
        submit(event){
            event.preventDefault();
            const data = {
                sender_username:username?.value,
                target_username:target?.value,
                body:file?.files[0].name
            }
            console.log("->",data)
            socket.send(JSON.stringify(data))
            },

        connect(){
            socket = new WebSocket(`ws://localhost:8080/chat?username=${username?.value}`);
            
            //listen for message
            socket.addEventListener('message',function(event){
                console.log('message from server', JSON.parse(event.data));
            });
            
            this.online++
            this.disabled = true
        }, 

        leave(){ 
            //esto es para ocultar ciertos componente a la hora de dejar el chat
            this.disabled = false
            this.online = 0 

            window.location.reload(true);    //actualización de la página
        },
    }
})

// Make connection
var socket = io.connect('http://10.253.69.155:4000');

// Query DOM
var message = document.getElementById('message'),
    user = document.getElementById('user'),
    btn = document.getElementById('send'),
    output = document.getElementById('output'),
    feedback = document.getElementById('feedback');

var joined = false;

function onSend(e) {
    if(user.value != '' && message.value != '') {
        if(joined == false) {
            joined = true;
            socket.emit('join',user.value);  
        }
           
        var date = new Date();
        var timestamp = Math.floor(date.getTime()/1000);
        
        const chatmsg = {
            user: user.value,
            timestamp: timestamp,
            message: message.value
        };

        const msgJSON = JSON.stringify(chatmsg);
        console.log(msgJSON);

        socket.emit('chat', msgJSON);
    }
    else {
        alert("Both Name and Message are required");
    }

    message.value = "";
}

// Button - to send a message
btn.addEventListener('click', onSend);

// message box
message.addEventListener('keypress', function(e){
    socket.emit('typing', name.value);

    if(e.keyCode==13) {
        onSend(e);
    }
})

// Listen for events 
socket.on('chat', function(data){
    feedback.innerHTML = '';

    var date = new Date(data.Timestamp * 1000);
    var timestr = date.getHours() + ':' + date.getMinutes() + ':' + date.getSeconds();
  
    if(data.EvtType == 'join')
        output.innerHTML += '<p>' + data.User + ' were joined   ' + '<strong>('+ timestr+')</strong></p>';
    else if(data.EvtType == 'leave')
        output.innerHTML += '<p>' + data.User + ' was left   ' + '<strong>('+ timestr+')</strong></p>';
    else if(data.EvtType == 'message')
        output.innerHTML += '<p><strong>'+ data.User  + ': </strong>' + data.Text +'     <strong>('+ timestr+')</strong></p>';       
  });

socket.on('typing', function(data){
    feedback.innerHTML = '<p><em>' + data + ' is typing a message...</em></p>';
});
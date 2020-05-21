// Make connection
var socket = io.connect('http://localhost:4000');

// Query DOM
var message = document.getElementById('message'),
    user = document.getElementById('user'),
    btn = document.getElementById('send'),
    output = document.getElementById('output'),
    feedback = document.getElementById('feedback'),
    title = document.getElementById('title');

var userID = localStorage.getItem('userName'); // set userID if exists 
if(userID != '')    {
    user.value = userID;

    socket.emit('join',user.value);  
}

function onSend(e) {
    if(user.value == '') {
        alert("Name is required");
    }
    else if(message.value != '') {
        if(userID != user.value) {
            socket.emit('join',user.value);  
            userID = user.value;  // update
            localStorage.setItem('userName', user.value)  // save userID
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
  
    if(data.EvtType == 'join' && data.User != '')
        output.innerHTML += '<p>' + data.User + ' was joined </p>';
    else if(data.EvtType == 'leave' && data.User != '')
        output.innerHTML += '<p>' + data.User + ' was left </p>';
    else if(data.EvtType == 'message') {
        if(data.User == userID)
            output.innerHTML += '<div style="display: flex; justify-content: flex-end"><p>'+ data.Text +'     <strong>('+ timestr+')</strong></p></div>';      
        else 
            output.innerHTML += '<p><strong>' + data.User  + ': </strong>' + data.Text +'     <strong>('+ timestr+')</strong></p>';      
    }
        
    output.scrollIntoView(false);
  });

// Listen for events 
socket.on('participant', function(data){
    title.textContent = 'Web Chat (' + data + ')';
    console.log('update participants');
});

socket.on('typing', function(data){
    feedback.innerHTML = '<p><em>' + data + ' is typing a message...</em></p>';
});
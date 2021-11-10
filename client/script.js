/*var element = document.getElementById("pFeed");*/
const gFeed = document.getElementById('gFeed');
const pFeed = document.getElementById('pFeed');
const uFeed = document.getElementById('uFeed');

const globalUrl = 'http://localhost:8080/global';
const personalUrl = 'http://localhost:8080/personal/b5bd831e-f799-49b1-a73a-dd85cee07f22';
const allUsersUrl = 'http://localhost:8080/allUsers';
const userUrl = 'http://localhost:8080/user/';
const createUserUrl = 'http://localhost:8080/user';
const newTweetUrl = 'http://localhost:8080/tweet/';
//const deleteTweetUrl = 'http://localhost:8080/tweet/';
const followUrl = 'http://localhost:8080/follow/';
const unfollowUrl = 'http://localhost:8080/unfollow/';

const personID = 'b5bd831e-f799-49b1-a73a-dd85cee07f22';

function createNode(element) {
    return document.createElement(element);
}

function createText(content) {
    return document.createTextNode(content);
}

function append(parent, el) {
    return parent.appendChild(el)
}

function reload() {
    
    window.location.reload();
}

function reloadAndHide() {
    let editDiv = document.getElementById("editScreen");
    editDiv.style.display = "none";
    window.location.reload();
}

function followUnfollow(status, id, friendID) {
    // idea: check if the element on page is follow or unfollow to decide what to do
    if (status == "follow") {
        var url = followUrl + id;

        var request = new Request(url, {
            method: 'POST',
            body: JSON.stringify({
                followID: friendID,
            }),
            headers: new Headers()
        });

        console.log(request);
        console.log(request.body);

        fetch(request)
        .then(function() {
            console.log("let's follow");
        })

        reload()

        // return new value of the button
        return "unfollow";
    }
    else {
        var url = unfollowUrl + id;

        var request = new Request(url, {
            method: 'POST',
            body: JSON.stringify({
                followID: friendID,
            }),
            headers: new Headers()
        });

        fetch(request)
        .then(function() {
            console.log("let's unfollow");
        })

        reload()

        return "follow";
    }
}


function editTweet(id, edits) {
    var url = newTweetUrl.concat(id);
    
    var request = new Request(url, {
        method: 'PUT',
        body: JSON.stringify({
            content: edits,
        }),
        headers: new Headers()
    });

    fetch(request)
    .then(function() {
        console.log("PLEASE");
    })
}


function editPopup(id, content) {
    // make an edit field popup on the page
    // div
    let editDiv = document.getElementById("editScreen");
    editDiv.style.display = "block";

    let editCard = document.createElement("div");
    editCard.setAttribute("class", "editCard");

    // header
    let editTitle = document.createElement("h6");
    //editTitle.setAttribute("value", "Edit Tweet");
    editTitle.innerHTML = "Edit Tweet";

    // form
    let editForm = document.createElement("form");
    editForm.setAttribute("class", "editForm");
    editForm.setAttribute("onsubmit", "reloadAndHide()");

    // textfield with pre filled spaces
    let editField = document.createElement("textarea");
    editField.innerHTML = content;
    editField.setAttribute("rows", "5");
    editField.setAttribute("cols", "50");
    editField.setAttribute("name", "theEdit");
    //editField.setAttribute("value", content);
    editField.setAttribute("class", "editField");

    // submit button
    let submitEdit = document.createElement("input");
    submitEdit.setAttribute("type", "submit");
    submitEdit.setAttribute("value", "OK");
    submitEdit.setAttribute("class", "submitEdit");

    // event listener
    // to submit the thingy
    submitEdit.addEventListener('click', function() {
        // Actual edittweet function here
        editTweet(id, editField.value);
    })

    // actually add all of this shit to the frontend
    append(editForm, editField);
    append(editForm, submitEdit);
    append(editCard, editTitle);
    append(editCard, editForm);

    var element = document.getElementById("editScreen");
    append(element, editCard);
}

// delete a tweet
function deleteTweet(id) {
    var url = newTweetUrl.concat(id);
    console.log(url);

    fetch(url, {
        method: 'DELETE',
      })
      .then(res => res.text()) // or res.json()
      .then(res => console.log(res))
}

// create a new tweet
function createTweet() {
    let inputID = document.getElementById("inputID").value;
    let inputContent = document.getElementById("inputContent").value;

    var url = newTweetUrl.concat(inputID);

    var request = new Request(url, {
        method: 'POST',
        body: JSON.stringify({
            content: inputContent,
        }),
        headers: new Headers()
    });

    fetch(request)
    .then(function() {
        console.log("please work");
    })
}

// create a new user
function createUser() {
    let inputName = document.getElementById("inputName").value;
    
    var request = new Request(createUserUrl, {
        method: 'POST',
        body: JSON.stringify({
            name: inputName,
        }),
        headers: new Headers()
    });

    console.log(request.body);

    fetch(request)
    .then(function() {
        console.log("please work");
    })
}


function userRoster() {
    /* call the API */
    fetch(allUsersUrl)
    .then((resp) => resp.json())
    .then(function(data) {
        let users = data.AllUsers;
        return users.map(async function(users) {
            let name = createText(users.name);
            let id = createText(users.userID);

            /* create contents */
            let userInfo = document.createElement("div");
            let div1 = document.createElement('div');
            let div2 = document.createElement('div');

            let follow = document.createElement('input');
            follow.setAttribute("type", "button");

            /* set whether to show follow or unfollow button here */
            var followStatusUrl = followUrl + personID + "/" + users.userID;

            /*
            var followRequest = new Request(followStatusUrl, {
                method: 'GET',
                body: JSON.stringify({
                    followID: users.userID,
                }),
                headers: new Headers()
            });
            */

            var followStatus = await fetch(followStatusUrl)
            .then((resp) => resp.json())
            .then(function(status) {
                console.log(status)
                return status.followed;
            })
            .catch(function(error) {
                console.log(error)
            })

            if (followStatus == "true") {
                follow.setAttribute("value", "unfollow");
            }
            else {
                follow.setAttribute("value", "follow");
            }

            //follow.setAttribute("value", "follow");
            follow.setAttribute("class", "followButton");

            follow.addEventListener('click', function() {
                // follow function here 
                returnValue = followUnfollow(follow.value, personID, users.userID);
                console.log(id);
                follow.setAttribute("value", returnValue);
            })

            userInfo.setAttribute("id", "userCard");
            div1.setAttribute("id", "userName");
            div2.setAttribute("id", "userID");
            
            /* add text node to divs */
            append(div1, name);
            append(div2, id);

            // add content to user card
            append(userInfo, div1);
            append(userInfo, div2);
            append(userInfo, follow);

            console.log(userInfo);

            // add card to the feed
            append(uFeed, userInfo);
        })
    })
}

function globalFeed() {
    /* call the API */
    fetch(globalUrl)
    /* get JSON data */
    .then((resp) => resp.json())
    .then(function(data) {
        /* do something with the data */
        let tweets = data.AllTweets;
        return tweets.map(async function(tweets) {
            /* get name of user */
            var nameOfUser = "";
            var userEndpoint = userUrl.concat(tweets.userID);
            console.log(userEndpoint);
            nameOfUser = await fetch(userEndpoint)
                /* fetch user info */
                .then((userResp) => userResp.json())
                .then(function(userData) {
                    console.log(userData)
                    console.log(userData.name)
                    return userData.name;
                })
                .catch(function(userError) {
                    console.log(userError);
                });

            let name = createText(nameOfUser);
            console.log(nameOfUser);
            /* get content of tweet */
            let content = createText(tweets.content);
            
            /* create elements */
            let card = document.createElement("div");
            let div1 = document.createElement('div');
            let div2 = document.createElement('div');

            let buttons = document.createElement('div');

            let editButton = document.createElement("input");
            editButton.setAttribute("type", "button");
            editButton.setAttribute("value", "edit");
            editButton.setAttribute("class", "editButton");

            editButton.addEventListener('click', function() {
                editPopup(tweets.tweetID, tweets.content);
            })

            let deleteForm = document.createElement("form");
            deleteForm.setAttribute("onsubmit", "reload()");
            deleteForm.setAttribute("class", "deleteForm");
            let deleteButton = document.createElement("input");
            deleteButton.setAttribute("type", "submit");
            deleteButton.setAttribute("value", "delete");
            deleteButton.setAttribute("class", "deleteButton");

            deleteButton.addEventListener('click', function() {
                deleteTweet(tweets.tweetID);
            })
            

            card.setAttribute("id", "tweetCard");
            div1.setAttribute("id", "tweetName");
            div2.setAttribute("id", "tweetContent");
            buttons.setAttribute("id", "crudButtons");
            
            /* add edit input button to buttons div */
            append(buttons, editButton);
            append(deleteForm, deleteButton);
            append(buttons, deleteForm);

            /* add text node to list element */
            append(div1, name);
            append(div2, content);

            /* add content to tweet card */
            append(card, div1);
            append(card, div2);
            append(card, buttons);

            /* add card to the feed */
            append(gFeed, card);
        })
    }) 
    .catch(function(error) {
        /* do something with the error */
        console.log(error);
    });
}

function personalFeed() {
    /* call the API */
    fetch(personalUrl)
    .then((resp) => resp.json())
    .then(function(data) {
        let tweets = data.AllTweets;
        return tweets.map(async function(tweets) {
            var nameOfUser = "";
            var userEndpoint = userUrl.concat(tweets.userID)
            console.log(userEndpoint);
            nameOfUser = await fetch(userEndpoint)
                .then((userResp) => userResp.json())
                .then(function(userData) {
                    return userData.name;
                })
                .catch(function(userError) {
                    console.log(userError);
                });

            let name = createText(nameOfUser);
            let content = createText(tweets.content);

            /* create elements */
            let card = document.createElement("div");
            let div1 = document.createElement('div');
            let div2 = document.createElement('div');

            let buttons = document.createElement('div');

            let editButton = document.createElement("input");
            editButton.setAttribute("type", "button");
            editButton.setAttribute("value", "edit");
            editButton.setAttribute("class", "editButton");

            editButton.addEventListener('click', function() {
                editPopup(tweets.tweetID, tweets.content);
            })

            let deleteForm = document.createElement("form");
            deleteForm.setAttribute("onsubmit", "reload()");
            deleteForm.setAttribute("class", "deleteForm");
            let deleteButton = document.createElement("input");
            deleteButton.setAttribute("type", "submit");
            deleteButton.setAttribute("value", "delete");
            deleteButton.setAttribute("class", "deleteButton");

            deleteButton.addEventListener('click', function() {
                deleteTweet(tweets.tweetID);
            })
            
            card.setAttribute("id", "tweetCard");
            div1.setAttribute("id", "tweetName");
            div2.setAttribute("id", "tweetContent");
            buttons.setAttribute("id", "crudButtons");

            /* add edit input button to buttons div */
            append(buttons, editButton);
            append(deleteForm, deleteButton);
            append(buttons, deleteForm);

            /* add text node to list element */
            append(div1, name);
            append(div2, content);

            /* add content to tweet card */
            append(card, div1);
            append(card, div2);
            append(card, buttons);

            /* add card to the feed */
            append(pFeed, card);
            
        })
    })
    .catch(function(error) {
        console.log(error);
    });
}

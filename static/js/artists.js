const cards = document.querySelectorAll('.artistCard');
const input = document.getElementById('search');
const boxSuggestion = document.querySelector('.boxSuggestion');
const artistParent = document.querySelector('.artist');
let boxArr = [];
let dataArtist = [];

const HandleCardClick = (card) => {
    const artistId = card.getAttribute("id");
    window.location.href = `/artist/${artistId}`;
};

const cardClickEvent = () => {
    const allCards = document.querySelectorAll('.artistCard');
    allCards.forEach(card => {
        card.addEventListener('click', () => HandleCardClick(card));
    });
};

const fetchData = async() => {
    const response = await fetch('/api/artists/');
    const resp = await response.json();
    dataArtist = resp;
};

const initializeData = async () => {
    await fetchData();
    cardClickEvent();
};

initializeData();

input.addEventListener('input', () => {
    const inputValue = input.value.toLowerCase();
    artistParent.innerHTML = '';  
    boxSuggestion.innerHTML = '';
    boxArr = []; 

    // Loop through dataArtist and filter by the input value
    dataArtist.forEach(obj => {
        const Name = obj.name.toLowerCase().includes(inputValue);
        const location = obj.Locations.toLowerCase();
        const date = obj.creationDate.toString().includes(inputValue);
        const firstAlbum = obj.firstAlbum.toLowerCase().includes(inputValue);
        const membersMatch = obj.members.some(member => member.toLowerCase().includes(inputValue));

        // Check if the artist's name matches the input
        if (Name && !boxArr.some(item => item.toLowerCase().includes(obj.name.toLowerCase()))) {
            boxArr.push(obj.name);
        }

        // Check locations
        const locationArr = obj.Locations.split(' ');
        locationArr.forEach(loc => {
            if (loc.toLowerCase().includes(inputValue) && !boxArr.some(item => item.toLowerCase().includes(loc.toLowerCase()))) {
                boxArr.push(loc);
            }
        });

        // Check creation date
        if (date && !boxArr.some(item => item.toLowerCase().includes(String(obj.creationDate)))) {
            boxArr.push(String(obj.creationDate));
        }

        // Check first album
        if (firstAlbum && !boxArr.some(item => item.toLowerCase().includes(obj.firstAlbum.toLowerCase()))) {
            boxArr.push(obj.firstAlbum);
        }

        // Check members
        obj.members.forEach(member => {
            if (member.toLowerCase().includes(inputValue) && !boxArr.some(item => item.toLowerCase().includes(member.toLowerCase()))) {
                boxArr.push(member);
            }
        });

        // If any condition matches, create the artist card
        if (Name || location.includes(inputValue) || date || firstAlbum || membersMatch) {
            artistParent.innerHTML += `
                <div class="artistCard" id="${obj.id}">
                    <img src="${obj.image}" alt="${obj.image}" loading="lazy">
                    <div class="name">${obj.name}</div>
                </div>
            `;
        }
    });

    // If input is empty, reset the boxArr
    if (inputValue.length === 0) {
        boxArr = [];
    }

    // Render suggestions
    if (boxArr.length > 0) {
        boxArr.forEach(item => {
            const suggestion = document.createElement('div');
            suggestion.textContent = item;
            suggestion.classList.add('suggestion');
            boxSuggestion.appendChild(suggestion);
        });
    }
    cardClickEvent();
});

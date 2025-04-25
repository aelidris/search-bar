const locationDiv = document.querySelector('#location')
const locationParent = document.querySelector('#locationParent')
const locationElement = locationDiv.textContent.replace(/[/[\]]/g,"").split(" ")
console.log(locationElement)
const dateDiv = document.querySelector('#date')
const dateParent = document.querySelector('#dateParent')
const dateElement = dateDiv.textContent.split(' ')

const relationDiv = document.querySelector('#relation')
const relationParent = document.querySelector('#relationParent')
const relationElement = relationDiv.textContent.split(' ')
const goBack = document.querySelector('.btn')

console.log(relationDiv.textContent)
const ulLocation = document.createElement('ul')
locationElement.forEach(location => {
    const liLocation = document.createElement('li')
    liLocation.textContent = location
    ulLocation.appendChild(liLocation)
})
locationDiv.textContent = ''
locationParent.appendChild(ulLocation)


const ulDate = document.createElement('ul')
dateElement.forEach(date => {
    const liDate = document.createElement('li')
    liDate.textContent = date.slice(1)
    ulDate.appendChild(liDate)
})
dateDiv.textContent = ''
dateParent.appendChild(ulDate)

let data
data = relationDiv.textContent.replace('map[', '').replace(/[\[]/g, '').slice(0, -2).split('] ')
console.log(data)
data = data.reduce((acc, locationANDdate) => {
    const [loc, dat] = locationANDdate.split(':')
    acc[loc] = dat
    return acc
}, {})
console.log(Object.entries(data))
for ([city, dateCity] of Object.entries(data)) {
    const ulRelation = document.createElement('ul')
    const div = document.createElement('div')
    div.classList.add('cityDate')
    const h4City = document.createElement('h4')
    h4City.classList.add('city')
    h4City.textContent = city
    dateCity.split(' ').forEach(item => {
        const liRelation = document.createElement('li')
        liRelation.textContent = item
        ulRelation.appendChild(liRelation)
    })
    div.append(h4City, ulRelation)
    relationParent.appendChild(div)
}
relationDiv.textContent = ''

goBack.addEventListener('click', () => {
    window.location.href = '/'
})
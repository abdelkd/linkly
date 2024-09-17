const regex = RegExp(/^https:\/\/[a-zA-Z0-9\-_]+(\.[a-zA-Z0-9\-_]+)+(\/[^\s]*)?\/?$/);

const urlInput = document.querySelector('.url-input')
const submitButton = document.querySelector('.submit-btn')
const serverResponse = document.querySelector('.server-response')
const errorText = document.querySelector('.server-error')
const copyLinkIcon = document.querySelector('.copy-link')
const copiedLinkIcon = document.querySelector('.copied-link')
const generatedLink = document.querySelector('.generated-link')

let urlText = ''


const toggleCopyIcons = () => {
    copyLinkIcon.classList.toggle('hidden')
    copiedLinkIcon.classList.toggle('hidden')
}

const onCopyLinkIconClick = () => {
  try {
    navigator.clipboard.writeText(generatedLink.textContent)
    toggleCopyIcons()
    setTimeout(toggleCopyIcons, 1500)

  } catch (err) {
    errorText.textContent = 'Could not copy URL.'
    console.error(err)
    return
  }
}

const handleUrlSubmit = async () => {
  // Reset states
  serverResponse.classList.add('hidden')
  serverResponse.classList.remove('flex')
  errorText.textContent = ''
  console.log(urlText, urlText.length, regex.test(urlText))


  try {
    if (urlText.length === 0) {
      throw new Error('URL cannot be empty.')
    }

    if (!regex.test(urlText)) {
      throw new Error('must be a valid URL.')
    }

    const response = await fetch('/v1/link', {
      method: 'POST',
      body: JSON.stringify({ location: urlText })
    }).then((r) => r.json())

    console.log({response})
    if (response?.location === null || response?.code === null) throw new Error('Bad Server Response')

    generatedLink.textContent = window.origin + '/' + response.code

    serverResponse.classList.add('flex')
    serverResponse.classList.remove('hidden')

  } catch (err) {
    errorText.textContent = err
    console.log("from catch", err)
    return
  }
}

submitButton.addEventListener('click', handleUrlSubmit)
copyLinkIcon.addEventListener('click', onCopyLinkIconClick)
urlInput.addEventListener('input', (e) => {
  urlText = e.target.value
})

const regex = RegExp(/^https:\/\/[a-zA-Z0-9\-_]+(\.[a-zA-Z0-9\-_]+)+(\/[^\s]*)?$/)

const urlInput = document.querySelector('.url-input')
const submitButton = document.querySelector('.submit-btn')
const serverResponse = document.querySelector('.server-response')
const errorText = document.querySelector('.server-error')
const copyLinkIcon = document.querySelector('.copy-link')
const copiedLinkIcon = document.querySelector('.copied-link')
const generatedLink = document.querySelector('.generated-link')

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


  const location = urlInput.value

  if (!regex.test(urlInput.value)) {
    errorText.textContent = 'URL cannot be empty.'
    console.error(errorText.textContent)
    return
  }


  try {
    const response = await fetch('/v1/link', {
      method: 'POST',
      body: JSON.stringify({ location })
    }).then((r) => r.json())

    if (response?.location === null || response?.code === null) throw new Error('Bad Server Response')

    generatedLink.textContent = window.origin + '/' + response.code

    serverResponse.classList.add('flex')
    serverResponse.classList.remove('hidden')

  } catch (err) {
    errorText.textContent = 'Could not create new link, try again.'
    console.error(err)
    return
  }
}

submitButton.addEventListener('click', handleUrlSubmit)
copyLinkIcon.addEventListener('click', onCopyLinkIconClick)

import svgReadyRaw from './inlineassets/clipboard-ready.svg'
import svgSuccessRaw from './inlineassets/clipboard-success.svg'

import ClipboardJS from 'clipboard'

(() => {
  if (!ClipboardJS.isSupported() || !document.querySelectorAll) {
    return
  }

  function svgdom(raw) {
    const div = document.createElement('div')
    div.innerHTML = raw
    div.firstChild.style.display = 'none'

    return div.firstChild
  }

  const svgReadyElement = svgdom(svgReadyRaw)
  const svgSuccessElement = svgdom(svgSuccessRaw)

  let indicatorDelay, indicatorReset = () => null

  const triggers = document.querySelectorAll('.snippet > button')

  Array.prototype.forEach.call(triggers, (button) => {
    const iconReady = svgReadyElement.cloneNode(true)
    const iconSuccess = svgSuccessElement.cloneNode(true)

    button.appendChild(iconReady)
    button.appendChild(iconSuccess)
    button.disabled = false
    iconReady.style.display = ''
  })

  const clipboard = new ClipboardJS(triggers, {
    text: (trigger) => trigger.parentNode.getElementsByTagName('input')[0].value,
  })
  clipboard.on('success', (e) => {
    e.clearSelection()

    const button = e.trigger
    const iconReady = button.childNodes[0]
    const iconSuccess = button.childNodes[1]

    iconReady.style.display = 'none'
    iconSuccess.style.display = ''

    indicatorReset()
    indicatorReset = () => {
      clearInterval(indicatorDelay)
      iconReady.style.display = ''
      iconSuccess.style.display = 'none'
      indicatorReset = () => null
    }

    indicatorDelay = setTimeout(indicatorReset, 2000)
  })
})()

import 'dart:io';
import 'dart:convert';

import 'package:flutter_inappwebview/flutter_inappwebview.dart';

import 'package:yordamchi/core/repository/repository.dart';
import 'constants.dart';
import 'crashlytics.dart';

final _loadedState = Platform.isIOS ? 1 : 2;
InAppWebViewController? _controller;
Object? _error;
int _state = 0;

Future<void> initSearch() async {
  final webView = HeadlessInAppWebView(
    initialSettings: InAppWebViewSettings(userAgent: _userAgent),
    onWebViewCreated: (controller) {
      _controller = controller;
    },
    onLoadStart: (controller, url) {
      _controller!.evaluateJavascript(source: _spoofSource);
    },
    onLoadStop: (controller, url) {
      _state++;
    },
    onReceivedError: (controller, request, error) {
      _error = error;
      _state++;
    },
  );
  await webView.run();
}

Future<String> getSearchResults(String query) async {
  _state = 0;
  _error = null;
  final encoded = Uri.encodeComponent(query);
  final url = 'https://www.google.com/search?q=$encoded&hl=${l10n.localeName}';
  await _controller!.loadUrl(urlRequest: URLRequest(url: WebUri(url)));
  while (_state < _loadedState && _error == null) {
    await Future.delayed(const Duration(milliseconds: 100));
  }
  if (_error != null) {
    Crashlytics.log(_error!, StackTrace.current);
    return '**No result**';
  }
  final json = await _controller!.evaluateJavascript(source: _searchSource);
  if (json == 'recaptcha') {
    Crashlytics.log('Recaptcha for the query: $query', StackTrace.current);
    return '**No result**';
  }
  var response = '';
  for (final result in jsonDecode(json)) {
    if (result['title'] != null) {
      response += '**Title**: ';
      response += result['title'];
      response += '\n';
    }
    if (result['content'] != null) {
      response += '**Content**: ';
      response += result['content'];
      response += '\n';
    }
    if (result['url'] != null) {
      response += '**URL**: ';
      response += result['url'];
      response += '\n';
    }
    response += '---\n';
  }
  return response;
}

Future<String> getPageSource() async {
  const pageSource = 'document.documentElement.outerHTML';
  return await _controller!.evaluateJavascript(source: pageSource);
}

final _regExp = RegExp(r'^Version (\d+)(?:\.(\d+)(?:\.(\d+))?)? \(Build .+\)$');

String get _userAgent {
  if (Platform.isIOS) {
    final match = _regExp.firstMatch(Platform.operatingSystemVersion);
    if (match != null) {
      final digit1 = match.group(1)!;
      final digit2 = match.group(2) ?? '0';
      final digit3 = match.group(3) ?? '0';
      return 'Mozilla/5.0 (iPhone; CPU iPhone OS ${digit1}_${digit2}_$digit3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/137.0.7151.107 Mobile/15E148 Safari/604.1';
    } else {
      return 'Mozilla/5.0 (iPhone; CPU iPhone OS 18_5_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/137.0.7151.107 Mobile/15E148 Safari/604.1';
    }
  } else {
    return 'Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Mobile Safari/537.36';
  }
}

final _languages = {'uz': 'uz-UZ', 'en': 'en-US', 'ru': 'ru-RU'};

String get _spoofSource =>
    '''
Object.defineProperty(navigator, 'webdriver', {
  get: () => undefined
});

Object.defineProperty(screen, 'pixelDepth', {
  get: () => 24
});

Object.defineProperty(screen, 'colorDepth', {
  get: () => 24
});

Object.defineProperty(screen, 'width', {
  get: () => ${Platform.isIOS ? '375' : '360'}
});

Object.defineProperty(screen, 'height', {
  get: () => ${Platform.isIOS ? '667' : '640'}
});

Object.defineProperty(screen, 'availWidth', {
  get: () => ${Platform.isIOS ? '375' : '360'}
});

Object.defineProperty(screen, 'availHeight', {
  get: () => ${Platform.isIOS ? '647' : '620'}
});

(function() {
  const getParameter = WebGLRenderingContext.prototype.getParameter;
  WebGLRenderingContext.prototype.getParameter = function(parameter) {
    if (parameter === 7936 /* GL_VENDOR */) {
      return ${Platform.isIOS ? "'Apple Inc.'" : "'Google Inc. (ANGLE)'"};
    }
    if (parameter === 7938 /* GL_VERSION */) {
      return 'WebGL 1.0 (OpenGL ES 2.0 Chromium)';
    }
    if (parameter === 35724 /* GL_SHADING_LANGUAGE_VERSION */) {
      return 'WebGL GLSL ES 1.0 (OpenGL ES GLSL ES 1.00 Chromium)';
    }
    if (parameter === 37445 /* UNMASKED_VENDOR_WEBGL */) {
      return ${Platform.isIOS ? "'Google Inc. (Apple)'" : "'Google Inc.'"};
    }
    if (parameter === 37446 /* UNMASKED_RENDERER_WEBGL */) {
      return ${Platform.isIOS ? "'ANGLE (Apple, ANGLE Metal Renderer: Apple M2, Unspecified Version)'" : "'ANGLE (Qualcomm, Adreno (TM) 650, OpenGL ES 3.2)'"};
    }
    return getParameter.call(this, parameter);
  };
})();

(function() {
  const originalGetContext = HTMLCanvasElement.prototype.getContext;
  HTMLCanvasElement.prototype.getContext = function(type, contextAttributes) {
    const context = originalGetContext.call(this, type, contextAttributes);
    if (type === '2d') {
      const originalGetImageData = context.getImageData;
      context.getImageData = function(sx, sy, sw, sh) {
        const imageData = originalGetImageData.call(this, sx, sy, sw, sh);
        for (let i = 0; i < imageData.data.length; i += 4) {
          imageData.data[i] += Math.floor((Math.random() - 0.5) * 2);
          imageData.data[i + 1] += Math.floor((Math.random() - 0.5) * 2);
          imageData.data[i + 2] += Math.floor((Math.random() - 0.5) * 2);
        }
        return imageData;
      };
    }
    return context;
  };
})();

(function() {
  const pluginsData = [
    {
      name: "PDF Viewer",
      filename: "internal-pdf-viewer",
      description: "Portable Document Format"
    },
    {
      name: "Chrome PDF Plugin",
      filename: "internal-pdf-viewer",
      description: "Portable Document Format"
    }
  ];

  function generatePluginArray() {
    const plugins = pluginsData.map(p => {
      const plugin = {};
      Object.setPrototypeOf(plugin, Plugin.prototype);
      Object.defineProperties(plugin, {
        name: { get: () => p.name },
        filename: { get: () => p.filename },
        description: { get: () => p.description },
        length: { get: () => 0 },
      });
      return plugin;
    });

    const pluginArray = {};
    Object.setPrototypeOf(pluginArray, PluginArray.prototype);
    Object.defineProperties(pluginArray, {
      length: { get: () => plugins.length },
      item: { value: (index) => plugins[index] },
      namedItem: { value: (name) => plugins.find(p => p.name === name) },
    });

    plugins.forEach((p, i) => pluginArray[i] = p);
    return pluginArray;
  }

  Object.defineProperty(navigator, 'plugins', {
    get: () => generatePluginArray()
  });

  Object.defineProperty(navigator, 'mimeTypes', {
    get: () => ({
      length: 2,
      'application/pdf': { type: 'application/pdf', suffixes: 'pdf', description: 'Portable Document Format', enabledPlugin: navigator.plugins[0] },
      'text/pdf': { type: 'text/pdf', suffixes: 'pdf', description: 'Portable Document Format', enabledPlugin: navigator.plugins[1] },
      item: function(index) { return this[Object.keys(this)[index]]; },
      namedItem: function(name) { return this[name]; }
    })
  });
})();

Object.defineProperty(navigator, 'platform', {
  get: () => ${Platform.isIOS ? "'iPhone'" : "'Linux armv81'"}
});

Object.defineProperty(navigator, 'hardwareConcurrency', {
  get: () => ${Platform.isIOS ? '6' : '8'}
});

Object.defineProperty(navigator, 'deviceMemory', {
  get: () => 4
});

Object.defineProperty(navigator, 'languages', {
  get: () => ['${_languages[l10n.localeName]}', '${l10n.localeName}']
});

Object.defineProperty(navigator, 'vendor', {
  get: () => 'Google Inc.'
});

Object.defineProperty(navigator, 'vendorSub', {
  get: () => ''
});

Object.defineProperty(navigator, 'productSub', {
  get: () => '20030107'
});

Object.defineProperty(Intl.DateTimeFormat.prototype, 'resolvedOptions', {
  value: function() {
    const options = Object.getOwnPropertyDescriptor(Intl.DateTimeFormat.prototype, 'resolvedOptions').value.call(this);
    options.timeZone = '$timezone';
    return options;
  }
});

Date.prototype.getTimezoneOffset = function() {
  return ${-DateTime.now().timeZoneOffset.inMinutes};
};

(function() {
  const audioContext = window.AudioContext || window.webkitAudioContext;
  if (audioContext) {
    const originalCreateAnalyser = audioContext.prototype.createAnalyser;
    audioContext.prototype.createAnalyser = function() {
      const analyser = originalCreateAnalyser.call(this);
      const originalGetFloatFrequencyData = analyser.getFloatFrequencyData;
      analyser.getFloatFrequencyData = function(array) {
        originalGetFloatFrequencyData.call(this, array);
        for (let i = 0; i < array.length; i++) {
          array[i] += (Math.random() - 0.5) * 0.0001;
        }
      };
      return analyser;
    };
  }
})();

Object.defineProperty(window.performance, 'now', {
  value: function() {
    return Date.now() + Math.random() * 10;
  }
});

Object.defineProperty(navigator, 'getBattery', {
  value: undefined
});

(function() {
  const originalOffsetWidth = Object.getOwnPropertyDescriptor(HTMLElement.prototype, 'offsetWidth');
  const originalOffsetHeight = Object.getOwnPropertyDescriptor(HTMLElement.prototype, 'offsetHeight');
  
  Object.defineProperty(HTMLElement.prototype, 'offsetWidth', {
    get: function() {
      const width = originalOffsetWidth.get.call(this);
      return width + (Math.random() - 0.5) * 0.1;
    }
  });
  
  Object.defineProperty(HTMLElement.prototype, 'offsetHeight', {
    get: function() {
      const height = originalOffsetHeight.get.call(this);
      return height + (Math.random() - 0.5) * 0.1;
    }
  });
})();

Object.defineProperty(navigator, 'connection', {
  get: () => ({
    effectiveType: '4g',
    downlink: 10,
    rtt: 50,
    saveData: false
  })
});

Object.defineProperty(navigator, 'permissions', {
  value: {
    query: () => Promise.resolve({ state: 'denied' })
  }
});

Object.defineProperty(navigator, 'geolocation', {
  value: {
    getCurrentPosition: function(success, error) {
      if (error) error({ code: 1, message: 'User denied Geolocation' });
    },
    watchPosition: function(success, error) {
      if (error) error({ code: 1, message: 'User denied Geolocation' });
    }
  }
});
''';

const _searchSource = r'''
(function() {
  var recaptcha = document.getElementById('recaptcha');
  if (recaptcha) {
    return 'recaptcha';  
  }

  function findElement(parent, tagName, className) {
    if (tagName && className) {
      var elems = parent.getElementsByTagName(tagName);
      for (var i = 0; i < elems.length; i++) {
        if (elems[i].className.includes(className)) {
          return elems[i];
        }
      }
    } else if (className) {
      var elems = parent.getElementsByClassName(className);
      return elems.length > 0 ? elems[0] : null;
    } else if (tagName) {
      var elems = parent.getElementsByTagName(tagName);
      return elems.length > 0 ? elems[0] : null;
    }
    return null;
  }

  function getTitle(element) {
    var elem = findElement(element, 'div', 'F0FGWb');
    return elem ? elem.innerText.trim() : '';
  }

  function getCleanContent(element) {
    var elem = findElement(element, 'div', 'VwiC3b');
    return elem ? elem.innerText.trim() : '';
  }

  function getContent(element) {
    return element.innerText.trim().replace(/\n{2,}/g, '\n');
  }

  function getUrl(element) {
    var elems = element.getElementsByTagName('a');
    var urls = new Set();
    for (var i = elems.length - 1; i >= 0; i--) {
      var url = elems[i].href;
      if (url &&
          !url.startsWith('https://www.google.com/search') &&
          !url.startsWith('https://search.app.goo.gl') &&
          !url.startsWith('https://maps.google.com/maps')) {
        urls.add(url);
        break;
      }
    }
    return urls.length == 0 ? '' : [...urls][0];
  }

  var results = [];
  var searchResults = document.getElementsByClassName('MjjYud');
  for (var i = 0; i < searchResults.length; i++) {
    var result = searchResults[i];
    var content = getCleanContent(result);
    if (content) {
      results.push({
        title: getTitle(result),
        content: content,
        url: getUrl(result)
      });
    } else {
      content = getContent(result);
      if (content &&
          !content.startsWith('Reklama') &&
          !content.startsWith('Boshqalar qidirgan') &&
          !content.startsWith('People also ask') &&
          !content.startsWith('Ads') &&
          !content.startsWith('People also search for') &&
          !content.startsWith('Вопросы по теме') &&
          !content.startsWith('Реклама') &&
          !content.startsWith('Другие также ищут')) {
        results.push({
          title: getTitle(result),
          content: content,
          url: getUrl(result)
        });
      }
    }
  }

  var cleanResults = [];
  for (var i = 0; i < results.length; i++) {
    var result = results[i];
    if (result.content) {
      if (!result.title) delete result.title;
      if (!result.url) delete result.url;
      cleanResults.push(result);
    }
  }

  return JSON.stringify(cleanResults);
})();''';

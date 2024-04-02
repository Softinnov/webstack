new EventSource('http://127.0.0.1:8888/esbuild').addEventListener('change', () => location.reload())

// ihm/node_modules/svelte/src/runtime/internal/utils.js
function noop() {
}
function run(fn) {
  return fn();
}
function blank_object() {
  return /* @__PURE__ */ Object.create(null);
}
function run_all(fns) {
  fns.forEach(run);
}
function is_function(thing) {
  return typeof thing === "function";
}
function safe_not_equal(a, b) {
  return a != a ? b == b : a !== b || a && typeof a === "object" || typeof a === "function";
}
function is_empty(obj) {
  return Object.keys(obj).length === 0;
}

// ihm/node_modules/svelte/src/runtime/internal/globals.js
var globals = typeof window !== "undefined" ? window : typeof globalThis !== "undefined" ? globalThis : (
  // @ts-ignore Node typings have this
  global
);

// ihm/node_modules/svelte/src/runtime/internal/ResizeObserverSingleton.js
var ResizeObserverSingleton = class _ResizeObserverSingleton {
  /**
   * @private
   * @readonly
   * @type {WeakMap<Element, import('./private.js').Listener>}
   */
  _listeners = "WeakMap" in globals ? /* @__PURE__ */ new WeakMap() : void 0;
  /**
   * @private
   * @type {ResizeObserver}
   */
  _observer = void 0;
  /** @type {ResizeObserverOptions} */
  options;
  /** @param {ResizeObserverOptions} options */
  constructor(options) {
    this.options = options;
  }
  /**
   * @param {Element} element
   * @param {import('./private.js').Listener} listener
   * @returns {() => void}
   */
  observe(element2, listener) {
    this._listeners.set(element2, listener);
    this._getObserver().observe(element2, this.options);
    return () => {
      this._listeners.delete(element2);
      this._observer.unobserve(element2);
    };
  }
  /**
   * @private
   */
  _getObserver() {
    return this._observer ?? (this._observer = new ResizeObserver((entries) => {
      for (const entry of entries) {
        _ResizeObserverSingleton.entries.set(entry.target, entry);
        this._listeners.get(entry.target)?.(entry);
      }
    }));
  }
};
ResizeObserverSingleton.entries = "WeakMap" in globals ? /* @__PURE__ */ new WeakMap() : void 0;

// ihm/node_modules/svelte/src/runtime/internal/dom.js
var is_hydrating = false;
function start_hydrating() {
  is_hydrating = true;
}
function end_hydrating() {
  is_hydrating = false;
}
function append(target, node) {
  target.appendChild(node);
}
function append_styles(target, style_sheet_id, styles) {
  const append_styles_to = get_root_for_style(target);
  if (!append_styles_to.getElementById(style_sheet_id)) {
    const style = element("style");
    style.id = style_sheet_id;
    style.textContent = styles;
    append_stylesheet(append_styles_to, style);
  }
}
function get_root_for_style(node) {
  if (!node)
    return document;
  const root = node.getRootNode ? node.getRootNode() : node.ownerDocument;
  if (root && /** @type {ShadowRoot} */
  root.host) {
    return (
      /** @type {ShadowRoot} */
      root
    );
  }
  return node.ownerDocument;
}
function append_stylesheet(node, style) {
  append(
    /** @type {Document} */
    node.head || node,
    style
  );
  return style.sheet;
}
function insert(target, node, anchor) {
  target.insertBefore(node, anchor || null);
}
function detach(node) {
  if (node.parentNode) {
    node.parentNode.removeChild(node);
  }
}
function destroy_each(iterations, detaching) {
  for (let i = 0; i < iterations.length; i += 1) {
    if (iterations[i])
      iterations[i].d(detaching);
  }
}
function element(name) {
  return document.createElement(name);
}
function text(data) {
  return document.createTextNode(data);
}
function space() {
  return text(" ");
}
function listen(node, event, handler, options) {
  node.addEventListener(event, handler, options);
  return () => node.removeEventListener(event, handler, options);
}
function attr(node, attribute, value) {
  if (value == null)
    node.removeAttribute(attribute);
  else if (node.getAttribute(attribute) !== value)
    node.setAttribute(attribute, value);
}
function children(element2) {
  return Array.from(element2.childNodes);
}
function set_data(text2, data) {
  data = "" + data;
  if (text2.data === data)
    return;
  text2.data = /** @type {string} */
  data;
}
function set_input_value(input, value) {
  input.value = value == null ? "" : value;
}
function set_style(node, key, value, important) {
  if (value == null) {
    node.style.removeProperty(key);
  } else {
    node.style.setProperty(key, value, important ? "important" : "");
  }
}
function get_custom_elements_slots(element2) {
  const result = {};
  element2.childNodes.forEach(
    /** @param {Element} node */
    (node) => {
      result[node.slot || "default"] = true;
    }
  );
  return result;
}

// ihm/node_modules/svelte/src/runtime/internal/lifecycle.js
var current_component;
function set_current_component(component) {
  current_component = component;
}

// ihm/node_modules/svelte/src/runtime/internal/scheduler.js
var dirty_components = [];
var binding_callbacks = [];
var render_callbacks = [];
var flush_callbacks = [];
var resolved_promise = /* @__PURE__ */ Promise.resolve();
var update_scheduled = false;
function schedule_update() {
  if (!update_scheduled) {
    update_scheduled = true;
    resolved_promise.then(flush);
  }
}
function add_render_callback(fn) {
  render_callbacks.push(fn);
}
var seen_callbacks = /* @__PURE__ */ new Set();
var flushidx = 0;
function flush() {
  if (flushidx !== 0) {
    return;
  }
  const saved_component = current_component;
  do {
    try {
      while (flushidx < dirty_components.length) {
        const component = dirty_components[flushidx];
        flushidx++;
        set_current_component(component);
        update(component.$$);
      }
    } catch (e) {
      dirty_components.length = 0;
      flushidx = 0;
      throw e;
    }
    set_current_component(null);
    dirty_components.length = 0;
    flushidx = 0;
    while (binding_callbacks.length)
      binding_callbacks.pop()();
    for (let i = 0; i < render_callbacks.length; i += 1) {
      const callback = render_callbacks[i];
      if (!seen_callbacks.has(callback)) {
        seen_callbacks.add(callback);
        callback();
      }
    }
    render_callbacks.length = 0;
  } while (dirty_components.length);
  while (flush_callbacks.length) {
    flush_callbacks.pop()();
  }
  update_scheduled = false;
  seen_callbacks.clear();
  set_current_component(saved_component);
}
function update($$) {
  if ($$.fragment !== null) {
    $$.update();
    run_all($$.before_update);
    const dirty = $$.dirty;
    $$.dirty = [-1];
    $$.fragment && $$.fragment.p($$.ctx, dirty);
    $$.after_update.forEach(add_render_callback);
  }
}
function flush_render_callbacks(fns) {
  const filtered = [];
  const targets = [];
  render_callbacks.forEach((c) => fns.indexOf(c) === -1 ? filtered.push(c) : targets.push(c));
  targets.forEach((c) => c());
  render_callbacks = filtered;
}

// ihm/node_modules/svelte/src/runtime/internal/transitions.js
var outroing = /* @__PURE__ */ new Set();
function transition_in(block, local) {
  if (block && block.i) {
    outroing.delete(block);
    block.i(local);
  }
}

// ihm/node_modules/svelte/src/runtime/internal/each.js
function ensure_array_like(array_like_or_iterator) {
  return array_like_or_iterator?.length !== void 0 ? array_like_or_iterator : Array.from(array_like_or_iterator);
}

// ihm/node_modules/svelte/src/shared/boolean_attributes.js
var _boolean_attributes = (
  /** @type {const} */
  [
    "allowfullscreen",
    "allowpaymentrequest",
    "async",
    "autofocus",
    "autoplay",
    "checked",
    "controls",
    "default",
    "defer",
    "disabled",
    "formnovalidate",
    "hidden",
    "inert",
    "ismap",
    "loop",
    "multiple",
    "muted",
    "nomodule",
    "novalidate",
    "open",
    "playsinline",
    "readonly",
    "required",
    "reversed",
    "selected"
  ]
);
var boolean_attributes = /* @__PURE__ */ new Set([..._boolean_attributes]);

// ihm/node_modules/svelte/src/runtime/internal/Component.js
function mount_component(component, target, anchor) {
  const { fragment, after_update } = component.$$;
  fragment && fragment.m(target, anchor);
  add_render_callback(() => {
    const new_on_destroy = component.$$.on_mount.map(run).filter(is_function);
    if (component.$$.on_destroy) {
      component.$$.on_destroy.push(...new_on_destroy);
    } else {
      run_all(new_on_destroy);
    }
    component.$$.on_mount = [];
  });
  after_update.forEach(add_render_callback);
}
function destroy_component(component, detaching) {
  const $$ = component.$$;
  if ($$.fragment !== null) {
    flush_render_callbacks($$.after_update);
    run_all($$.on_destroy);
    $$.fragment && $$.fragment.d(detaching);
    $$.on_destroy = $$.fragment = null;
    $$.ctx = [];
  }
}
function make_dirty(component, i) {
  if (component.$$.dirty[0] === -1) {
    dirty_components.push(component);
    schedule_update();
    component.$$.dirty.fill(0);
  }
  component.$$.dirty[i / 31 | 0] |= 1 << i % 31;
}
function init(component, options, instance2, create_fragment2, not_equal, props, append_styles2 = null, dirty = [-1]) {
  const parent_component = current_component;
  set_current_component(component);
  const $$ = component.$$ = {
    fragment: null,
    ctx: [],
    // state
    props,
    update: noop,
    not_equal,
    bound: blank_object(),
    // lifecycle
    on_mount: [],
    on_destroy: [],
    on_disconnect: [],
    before_update: [],
    after_update: [],
    context: new Map(options.context || (parent_component ? parent_component.$$.context : [])),
    // everything else
    callbacks: blank_object(),
    dirty,
    skip_bound: false,
    root: options.target || parent_component.$$.root
  };
  append_styles2 && append_styles2($$.root);
  let ready = false;
  $$.ctx = instance2 ? instance2(component, options.props || {}, (i, ret, ...rest) => {
    const value = rest.length ? rest[0] : ret;
    if ($$.ctx && not_equal($$.ctx[i], $$.ctx[i] = value)) {
      if (!$$.skip_bound && $$.bound[i])
        $$.bound[i](value);
      if (ready)
        make_dirty(component, i);
    }
    return ret;
  }) : [];
  $$.update();
  ready = true;
  run_all($$.before_update);
  $$.fragment = create_fragment2 ? create_fragment2($$.ctx) : false;
  if (options.target) {
    if (options.hydrate) {
      start_hydrating();
      const nodes = children(options.target);
      $$.fragment && $$.fragment.l(nodes);
      nodes.forEach(detach);
    } else {
      $$.fragment && $$.fragment.c();
    }
    if (options.intro)
      transition_in(component.$$.fragment);
    mount_component(component, options.target, options.anchor);
    end_hydrating();
    flush();
  }
  set_current_component(parent_component);
}
var SvelteElement;
if (typeof HTMLElement === "function") {
  SvelteElement = class extends HTMLElement {
    /** The Svelte component constructor */
    $$ctor;
    /** Slots */
    $$s;
    /** The Svelte component instance */
    $$c;
    /** Whether or not the custom element is connected */
    $$cn = false;
    /** Component props data */
    $$d = {};
    /** `true` if currently in the process of reflecting component props back to attributes */
    $$r = false;
    /** @type {Record<string, CustomElementPropDefinition>} Props definition (name, reflected, type etc) */
    $$p_d = {};
    /** @type {Record<string, Function[]>} Event listeners */
    $$l = {};
    /** @type {Map<Function, Function>} Event listener unsubscribe functions */
    $$l_u = /* @__PURE__ */ new Map();
    constructor($$componentCtor, $$slots, use_shadow_dom) {
      super();
      this.$$ctor = $$componentCtor;
      this.$$s = $$slots;
      if (use_shadow_dom) {
        this.attachShadow({ mode: "open" });
      }
    }
    addEventListener(type, listener, options) {
      this.$$l[type] = this.$$l[type] || [];
      this.$$l[type].push(listener);
      if (this.$$c) {
        const unsub = this.$$c.$on(type, listener);
        this.$$l_u.set(listener, unsub);
      }
      super.addEventListener(type, listener, options);
    }
    removeEventListener(type, listener, options) {
      super.removeEventListener(type, listener, options);
      if (this.$$c) {
        const unsub = this.$$l_u.get(listener);
        if (unsub) {
          unsub();
          this.$$l_u.delete(listener);
        }
      }
    }
    async connectedCallback() {
      this.$$cn = true;
      if (!this.$$c) {
        let create_slot = function(name) {
          return () => {
            let node;
            const obj = {
              c: function create() {
                node = element("slot");
                if (name !== "default") {
                  attr(node, "name", name);
                }
              },
              /**
               * @param {HTMLElement} target
               * @param {HTMLElement} [anchor]
               */
              m: function mount(target, anchor) {
                insert(target, node, anchor);
              },
              d: function destroy(detaching) {
                if (detaching) {
                  detach(node);
                }
              }
            };
            return obj;
          };
        };
        await Promise.resolve();
        if (!this.$$cn || this.$$c) {
          return;
        }
        const $$slots = {};
        const existing_slots = get_custom_elements_slots(this);
        for (const name of this.$$s) {
          if (name in existing_slots) {
            $$slots[name] = [create_slot(name)];
          }
        }
        for (const attribute of this.attributes) {
          const name = this.$$g_p(attribute.name);
          if (!(name in this.$$d)) {
            this.$$d[name] = get_custom_element_value(name, attribute.value, this.$$p_d, "toProp");
          }
        }
        for (const key in this.$$p_d) {
          if (!(key in this.$$d) && this[key] !== void 0) {
            this.$$d[key] = this[key];
            delete this[key];
          }
        }
        this.$$c = new this.$$ctor({
          target: this.shadowRoot || this,
          props: {
            ...this.$$d,
            $$slots,
            $$scope: {
              ctx: []
            }
          }
        });
        const reflect_attributes = () => {
          this.$$r = true;
          for (const key in this.$$p_d) {
            this.$$d[key] = this.$$c.$$.ctx[this.$$c.$$.props[key]];
            if (this.$$p_d[key].reflect) {
              const attribute_value = get_custom_element_value(
                key,
                this.$$d[key],
                this.$$p_d,
                "toAttribute"
              );
              if (attribute_value == null) {
                this.removeAttribute(this.$$p_d[key].attribute || key);
              } else {
                this.setAttribute(this.$$p_d[key].attribute || key, attribute_value);
              }
            }
          }
          this.$$r = false;
        };
        this.$$c.$$.after_update.push(reflect_attributes);
        reflect_attributes();
        for (const type in this.$$l) {
          for (const listener of this.$$l[type]) {
            const unsub = this.$$c.$on(type, listener);
            this.$$l_u.set(listener, unsub);
          }
        }
        this.$$l = {};
      }
    }
    // We don't need this when working within Svelte code, but for compatibility of people using this outside of Svelte
    // and setting attributes through setAttribute etc, this is helpful
    attributeChangedCallback(attr2, _oldValue, newValue) {
      if (this.$$r)
        return;
      attr2 = this.$$g_p(attr2);
      this.$$d[attr2] = get_custom_element_value(attr2, newValue, this.$$p_d, "toProp");
      this.$$c?.$set({ [attr2]: this.$$d[attr2] });
    }
    disconnectedCallback() {
      this.$$cn = false;
      Promise.resolve().then(() => {
        if (!this.$$cn) {
          this.$$c.$destroy();
          this.$$c = void 0;
        }
      });
    }
    $$g_p(attribute_name) {
      return Object.keys(this.$$p_d).find(
        (key) => this.$$p_d[key].attribute === attribute_name || !this.$$p_d[key].attribute && key.toLowerCase() === attribute_name
      ) || attribute_name;
    }
  };
}
function get_custom_element_value(prop, value, props_definition, transform) {
  const type = props_definition[prop]?.type;
  value = type === "Boolean" && typeof value !== "boolean" ? value != null : value;
  if (!transform || !props_definition[prop]) {
    return value;
  } else if (transform === "toAttribute") {
    switch (type) {
      case "Object":
      case "Array":
        return value == null ? null : JSON.stringify(value);
      case "Boolean":
        return value ? "" : null;
      case "Number":
        return value == null ? null : value;
      default:
        return value;
    }
  } else {
    switch (type) {
      case "Object":
      case "Array":
        return value && JSON.parse(value);
      case "Boolean":
        return value;
      case "Number":
        return value != null ? +value : value;
      default:
        return value;
    }
  }
}
function create_custom_element(Component, props_definition, slots, accessors, use_shadow_dom, extend) {
  let Class = class extends SvelteElement {
    constructor() {
      super(Component, slots, use_shadow_dom);
      this.$$p_d = props_definition;
    }
    static get observedAttributes() {
      return Object.keys(props_definition).map(
        (key) => (props_definition[key].attribute || key).toLowerCase()
      );
    }
  };
  Object.keys(props_definition).forEach((prop) => {
    Object.defineProperty(Class.prototype, prop, {
      get() {
        return this.$$c && prop in this.$$c ? this.$$c[prop] : this.$$d[prop];
      },
      set(value) {
        value = get_custom_element_value(prop, value, props_definition);
        this.$$d[prop] = value;
        this.$$c?.$set({ [prop]: value });
      }
    });
  });
  accessors.forEach((accessor) => {
    Object.defineProperty(Class.prototype, accessor, {
      get() {
        return this.$$c?.[accessor];
      }
    });
  });
  if (extend) {
    Class = extend(Class);
  }
  Component.element = /** @type {any} */
  Class;
  return Class;
}
var SvelteComponent = class {
  /**
   * ### PRIVATE API
   *
   * Do not use, may change at any time
   *
   * @type {any}
   */
  $$ = void 0;
  /**
   * ### PRIVATE API
   *
   * Do not use, may change at any time
   *
   * @type {any}
   */
  $$set = void 0;
  /** @returns {void} */
  $destroy() {
    destroy_component(this, 1);
    this.$destroy = noop;
  }
  /**
   * @template {Extract<keyof Events, string>} K
   * @param {K} type
   * @param {((e: Events[K]) => void) | null | undefined} callback
   * @returns {() => void}
   */
  $on(type, callback) {
    if (!is_function(callback)) {
      return noop;
    }
    const callbacks = this.$$.callbacks[type] || (this.$$.callbacks[type] = []);
    callbacks.push(callback);
    return () => {
      const index = callbacks.indexOf(callback);
      if (index !== -1)
        callbacks.splice(index, 1);
    };
  }
  /**
   * @param {Partial<Props>} props
   * @returns {void}
   */
  $set(props) {
    if (this.$$set && !is_empty(props)) {
      this.$$.skip_bound = true;
      this.$$set(props);
      this.$$.skip_bound = false;
    }
  }
};

// ihm/node_modules/svelte/src/shared/version.js
var PUBLIC_VERSION = "4";

// ihm/node_modules/svelte/src/runtime/internal/disclose-version/index.js
if (typeof window !== "undefined")
  (window.__svelte || (window.__svelte = { v: /* @__PURE__ */ new Set() })).v.add(PUBLIC_VERSION);

// ihm/app.svelte
function add_css(target) {
  append_styles(target, "svelte-1ur7k5m", 'h1.svelte-1ur7k5m.svelte-1ur7k5m{font-size:70px}p.svelte-1ur7k5m.svelte-1ur7k5m{font-size:large;left:45%;bottom:5%;position:fixed}.button.svelte-1ur7k5m.svelte-1ur7k5m{height:35px;width:35px;border:none;border-radius:5%;background-color:white;margin-left:3%}.button.svelte-1ur7k5m.svelte-1ur7k5m:hover{background-color:rgba(146, 146, 146, 0.381);cursor:pointer}.centered.svelte-1ur7k5m.svelte-1ur7k5m{width:30em;margin:auto;display:grid}.newTask.svelte-1ur7k5m.svelte-1ur7k5m{flex:0.75;left:12%;margin-bottom:15%;margin-right:1%;position:relative}.todos.svelte-1ur7k5m.svelte-1ur7k5m{margin-top:10%}.ajout.svelte-1ur7k5m.svelte-1ur7k5m{height:35px;width:35px;border:none;border-radius:5%;background-color:white;margin-left:2%;position:relative}.ajout.svelte-1ur7k5m.svelte-1ur7k5m:hover{background-color:rgba(146, 146, 146, 0.381);cursor:pointer}.ajout.svelte-1ur7k5m.svelte-1ur7k5m:disabled{background-color:white;color:rgba(128, 128, 128, 0.836);cursor:default}ul.svelte-1ur7k5m.svelte-1ur7k5m{max-height:15em;overflow:hidden;overflow-y:visible}li.svelte-1ur7k5m.svelte-1ur7k5m{display:flex;margin:3%}input[type="text"].svelte-1ur7k5m.svelte-1ur7k5m{flex:0.75;padding:0.5em;margin:-0.2em 0;border:none;font-size:large}.priority-btn.svelte-1ur7k5m.svelte-1ur7k5m{position:relative;width:10px;height:10px;border-radius:50%;margin:3%}.priority-btn.svelte-1ur7k5m.svelte-1ur7k5m:hover{cursor:pointer}.priority.svelte-1ur7k5m.svelte-1ur7k5m{display:inline;margin-left:13%;justify-content:center}.priority.svelte-1ur7k5m button.svelte-1ur7k5m{position:relative;width:10px;height:10px;border-radius:50%;margin-right:1%;margin-left:1%;bottom:10%}.priority.svelte-1ur7k5m button.svelte-1ur7k5m:hover{cursor:pointer}.urgent.svelte-1ur7k5m.svelte-1ur7k5m{background-color:rgba(255, 0, 0, 0.475)}.prioritaire.svelte-1ur7k5m.svelte-1ur7k5m{background-color:rgba(255, 255, 0, 0.475)}.nonprioritaire.svelte-1ur7k5m.svelte-1ur7k5m{background-color:rgba(0, 128, 0, 0.475)}.selectedurg.svelte-1ur7k5m.svelte-1ur7k5m{border:2px solid;border-color:black;background-color:red;transform:scale(1.15)}.selectedprio.svelte-1ur7k5m.svelte-1ur7k5m{border:2px solid;border-color:black;background-color:yellow;transform:scale(1.15)}.selectednonurg.svelte-1ur7k5m.svelte-1ur7k5m{border:2px solid;border-color:black;background-color:green;transform:scale(1.15)}');
}
function get_each_context(ctx, list, i) {
  const child_ctx = ctx.slice();
  child_ctx[20] = list[i];
  child_ctx[21] = list;
  child_ctx[22] = i;
  return child_ctx;
}
function create_each_block(ctx) {
  let li;
  let input;
  let t0;
  let button0;
  let t1;
  let button1;
  let t3;
  let button2;
  let t5;
  let mounted;
  let dispose;
  function input_input_handler_1() {
    ctx[15].call(
      input,
      /*each_value*/
      ctx[21],
      /*item_index*/
      ctx[22]
    );
  }
  function keydown_handler_1(...args) {
    return (
      /*keydown_handler_1*/
      ctx[16](
        /*item*/
        ctx[20],
        ...args
      )
    );
  }
  function click_handler_3() {
    return (
      /*click_handler_3*/
      ctx[17](
        /*item*/
        ctx[20]
      )
    );
  }
  return {
    c() {
      li = element("li");
      input = element("input");
      t0 = space();
      button0 = element("button");
      t1 = space();
      button1 = element("button");
      button1.textContent = "\u270F\uFE0F";
      t3 = space();
      button2 = element("button");
      button2.textContent = "\u{1F5D1}\uFE0F";
      t5 = space();
      attr(input, "id", "todo");
      attr(input, "type", "text");
      attr(input, "class", "svelte-1ur7k5m");
      attr(button0, "class", "priority-btn svelte-1ur7k5m");
      set_style(button0, "background-color", getPriorityColor(
        /*item*/
        ctx[20].priority
      ));
      attr(button1, "class", "button svelte-1ur7k5m");
      attr(button2, "class", "button svelte-1ur7k5m");
      attr(li, "class", "svelte-1ur7k5m");
    },
    m(target, anchor) {
      insert(target, li, anchor);
      append(li, input);
      set_input_value(
        input,
        /*item*/
        ctx[20].text
      );
      append(li, t0);
      append(li, button0);
      append(li, t1);
      append(li, button1);
      append(li, t3);
      append(li, button2);
      append(li, t5);
      if (!mounted) {
        dispose = [
          listen(input, "input", input_input_handler_1),
          listen(input, "keydown", keydown_handler_1),
          listen(button0, "click", click_handler_3),
          listen(button1, "click", function() {
            if (is_function(
              /*modify*/
              ctx[8](
                /*item*/
                ctx[20]
              )
            ))
              ctx[8](
                /*item*/
                ctx[20]
              ).apply(this, arguments);
          }),
          listen(button2, "click", function() {
            if (is_function(
              /*xclear*/
              ctx[7](
                /*item*/
                ctx[20]
              )
            ))
              ctx[7](
                /*item*/
                ctx[20]
              ).apply(this, arguments);
          })
        ];
        mounted = true;
      }
    },
    p(new_ctx, dirty) {
      ctx = new_ctx;
      if (dirty & /*todos*/
      1 && input.value !== /*item*/
      ctx[20].text) {
        set_input_value(
          input,
          /*item*/
          ctx[20].text
        );
      }
      if (dirty & /*todos*/
      1) {
        set_style(button0, "background-color", getPriorityColor(
          /*item*/
          ctx[20].priority
        ));
      }
    },
    d(detaching) {
      if (detaching) {
        detach(li);
      }
      mounted = false;
      run_all(dispose);
    }
  };
}
function create_fragment(ctx) {
  let div2;
  let h1;
  let t1;
  let div1;
  let input;
  let t2;
  let div0;
  let button0;
  let button0_class_value;
  let t3;
  let button1;
  let button1_class_value;
  let t4;
  let button2;
  let button2_class_value;
  let t5;
  let button3;
  let t6;
  let button3_disabled_value;
  let t7;
  let ul;
  let t8;
  let p;
  let t9;
  let t10;
  let mounted;
  let dispose;
  let each_value = ensure_array_like(
    /*todos*/
    ctx[0]
  );
  let each_blocks = [];
  for (let i = 0; i < each_value.length; i += 1) {
    each_blocks[i] = create_each_block(get_each_context(ctx, each_value, i));
  }
  return {
    c() {
      div2 = element("div");
      h1 = element("h1");
      h1.textContent = "My TodoList";
      t1 = space();
      div1 = element("div");
      input = element("input");
      t2 = space();
      div0 = element("div");
      button0 = element("button");
      t3 = space();
      button1 = element("button");
      t4 = space();
      button2 = element("button");
      t5 = space();
      button3 = element("button");
      t6 = text("\u2714\uFE0F");
      t7 = space();
      ul = element("ul");
      for (let i = 0; i < each_blocks.length; i += 1) {
        each_blocks[i].c();
      }
      t8 = space();
      p = element("p");
      t9 = text(
        /*remaining*/
        ctx[3]
      );
      t10 = text(" t\xE2ches restantes !");
      attr(h1, "class", "svelte-1ur7k5m");
      attr(input, "class", "newTask svelte-1ur7k5m");
      attr(input, "type", "text");
      attr(input, "placeholder", "Quoi d'autre?");
      attr(button0, "class", button0_class_value = "urgent " + /*selectedPriority*/
      (ctx[2] === 3 ? "selectedurg" : "") + " svelte-1ur7k5m");
      attr(button1, "class", button1_class_value = "prioritaire " + /*selectedPriority*/
      (ctx[2] === 2 ? "selectedprio" : "") + " svelte-1ur7k5m");
      attr(button2, "class", button2_class_value = "nonprioritaire " + /*selectedPriority*/
      (ctx[2] === 1 ? "selectednonurg" : "") + " svelte-1ur7k5m");
      attr(div0, "class", "priority svelte-1ur7k5m");
      attr(button3, "class", "ajout svelte-1ur7k5m");
      button3.disabled = button3_disabled_value = /*nouvelleTache*/
      ctx[1] == "";
      attr(ul, "id", "todo-list");
      attr(ul, "class", "todos svelte-1ur7k5m");
      attr(p, "class", "svelte-1ur7k5m");
      attr(div2, "class", "centered svelte-1ur7k5m");
    },
    m(target, anchor) {
      insert(target, div2, anchor);
      append(div2, h1);
      append(div2, t1);
      append(div2, div1);
      append(div1, input);
      set_input_value(
        input,
        /*nouvelleTache*/
        ctx[1]
      );
      append(div1, t2);
      append(div1, div0);
      append(div0, button0);
      append(div0, t3);
      append(div0, button1);
      append(div0, t4);
      append(div0, button2);
      append(div1, t5);
      append(div1, button3);
      append(button3, t6);
      append(div2, t7);
      append(div2, ul);
      for (let i = 0; i < each_blocks.length; i += 1) {
        if (each_blocks[i]) {
          each_blocks[i].m(ul, null);
        }
      }
      append(div2, t8);
      append(div2, p);
      append(p, t9);
      append(p, t10);
      if (!mounted) {
        dispose = [
          listen(
            input,
            "input",
            /*input_input_handler*/
            ctx[10]
          ),
          listen(
            input,
            "keydown",
            /*keydown_handler*/
            ctx[11]
          ),
          listen(
            button0,
            "click",
            /*click_handler*/
            ctx[12]
          ),
          listen(
            button1,
            "click",
            /*click_handler_1*/
            ctx[13]
          ),
          listen(
            button2,
            "click",
            /*click_handler_2*/
            ctx[14]
          ),
          listen(
            button3,
            "click",
            /*add*/
            ctx[6]
          )
        ];
        mounted = true;
      }
    },
    p(ctx2, [dirty]) {
      if (dirty & /*nouvelleTache*/
      2 && input.value !== /*nouvelleTache*/
      ctx2[1]) {
        set_input_value(
          input,
          /*nouvelleTache*/
          ctx2[1]
        );
      }
      if (dirty & /*selectedPriority*/
      4 && button0_class_value !== (button0_class_value = "urgent " + /*selectedPriority*/
      (ctx2[2] === 3 ? "selectedurg" : "") + " svelte-1ur7k5m")) {
        attr(button0, "class", button0_class_value);
      }
      if (dirty & /*selectedPriority*/
      4 && button1_class_value !== (button1_class_value = "prioritaire " + /*selectedPriority*/
      (ctx2[2] === 2 ? "selectedprio" : "") + " svelte-1ur7k5m")) {
        attr(button1, "class", button1_class_value);
      }
      if (dirty & /*selectedPriority*/
      4 && button2_class_value !== (button2_class_value = "nonprioritaire " + /*selectedPriority*/
      (ctx2[2] === 1 ? "selectednonurg" : "") + " svelte-1ur7k5m")) {
        attr(button2, "class", button2_class_value);
      }
      if (dirty & /*nouvelleTache*/
      2 && button3_disabled_value !== (button3_disabled_value = /*nouvelleTache*/
      ctx2[1] == "")) {
        button3.disabled = button3_disabled_value;
      }
      if (dirty & /*xclear, todos, modify, getPriorityColor, changePriority, handleKeydown*/
      929) {
        each_value = ensure_array_like(
          /*todos*/
          ctx2[0]
        );
        let i;
        for (i = 0; i < each_value.length; i += 1) {
          const child_ctx = get_each_context(ctx2, each_value, i);
          if (each_blocks[i]) {
            each_blocks[i].p(child_ctx, dirty);
          } else {
            each_blocks[i] = create_each_block(child_ctx);
            each_blocks[i].c();
            each_blocks[i].m(ul, null);
          }
        }
        for (; i < each_blocks.length; i += 1) {
          each_blocks[i].d(1);
        }
        each_blocks.length = each_value.length;
      }
      if (dirty & /*remaining*/
      8)
        set_data(
          t9,
          /*remaining*/
          ctx2[3]
        );
    },
    i: noop,
    o: noop,
    d(detaching) {
      if (detaching) {
        detach(div2);
      }
      destroy_each(each_blocks, detaching);
      mounted = false;
      run_all(dispose);
    }
  };
}
function getPriorityColor(priority) {
  switch (priority) {
    case 3:
      return "rgba(255, 0, 0)";
    case 2:
      return "rgba(255, 255, 0)";
    case 1:
      return "rgba(0, 128, 0)";
    default:
      return "rgba(0, 0, 0)";
  }
}
function customQueryEscape(str) {
  const uriStr = encodeURIComponent(str);
  const finalQueryStr = uriStr.replace(/!/g, "%21").replace(/'/g, "%27").replace(/\(/g, "%28").replace(/\)/g, "%29").replace(/\*/g, "%2A").replace(/%20/g, "+");
  return finalQueryStr;
}
function instance($$self, $$props, $$invalidate) {
  let remaining;
  let todos = [];
  let nouvelleTache = "";
  let selectedPriority = 2;
  function selectPriority(priority) {
    $$invalidate(2, selectedPriority = priority);
  }
  function changePriority(item) {
    if (item.priority >= 3) {
      item.priority = 0;
      item.priority += 1;
    } else {
      item.priority += 1;
    }
    console.log(item.priority);
    modify(item);
  }
  function add() {
    let todo = {
      text: nouvelleTache,
      priority: selectedPriority
    };
    try {
      sendTodo(todo, "add");
    } catch (error) {
      console.error(`Erreur lors de la connection au serveur : ${error.message}`);
    }
    $$invalidate(1, nouvelleTache = "");
  }
  function xclear(item) {
    try {
      sendTodo(item, "delete");
    } catch (error) {
      console.error(`Erreur lors de la connection au serveur : ${error.message}`);
    }
  }
  function modify(item) {
    try {
      sendTodo(item, "modify");
    } catch (error) {
      console.error(`Erreur lors de la connection au serveur : ${error.message}`);
    }
  }
  async function getTodoList() {
    const url = `/todos`;
    try {
      const reponse = await fetch(url, { method: "GET" });
      if (!reponse.ok) {
        const errorData = await reponse.text();
        alert(errorData);
        throw new Error(`Erreur lors de la requ\xEAte : ${reponse.status} ${reponse.statusText}`);
      }
      const result = await reponse.json();
      if (result != null) {
        $$invalidate(0, todos = result);
      }
      if (result == null) {
        $$invalidate(0, todos = []);
      }
    } catch (error) {
      console.error(`Une erreur s'est produite : ${error.message}`);
    }
  }
  async function sendTodo(todo, route) {
    todo.text = customQueryEscape(todo.text);
    if (route == "add") {
      var url = `/${route}?text=${todo.text}&priority=${todo.priority}`;
    } else {
      url = `/${route}?id=${todo.id}&text=${todo.text}&priority=${todo.priority}`;
    }
    try {
      const reponse = await fetch(url, { method: "POST" });
      if (!reponse.ok) {
        const errorData = await reponse.text();
        alert(errorData);
        throw new Error(`Erreur lors de la requ\xEAte : ${reponse.status} ${reponse.statusText}`);
      }
      await reponse.json();
      getTodoList();
    } catch (error) {
      console.error(`Une erreur s'est produite : ${error.message}`);
    }
  }
  function handleKeydown(e, item) {
    if (e.key == "Enter") {
      if (item != null) {
        modify(item);
      } else {
        add();
      }
    }
  }
  getTodoList();
  function input_input_handler() {
    nouvelleTache = this.value;
    $$invalidate(1, nouvelleTache);
  }
  const keydown_handler = (e) => handleKeydown(e, null);
  const click_handler = () => selectPriority(3);
  const click_handler_1 = () => selectPriority(2);
  const click_handler_2 = () => selectPriority(1);
  function input_input_handler_1(each_value, item_index) {
    each_value[item_index].text = this.value;
    $$invalidate(0, todos);
  }
  const keydown_handler_1 = (item, e) => handleKeydown(e, item);
  const click_handler_3 = (item) => changePriority(item);
  $$self.$$.update = () => {
    if ($$self.$$.dirty & /*todos*/
    1) {
      $:
        $$invalidate(3, remaining = todos.length);
    }
  };
  return [
    todos,
    nouvelleTache,
    selectedPriority,
    remaining,
    selectPriority,
    changePriority,
    add,
    xclear,
    modify,
    handleKeydown,
    input_input_handler,
    keydown_handler,
    click_handler,
    click_handler_1,
    click_handler_2,
    input_input_handler_1,
    keydown_handler_1,
    click_handler_3
  ];
}
var App = class extends SvelteComponent {
  constructor(options) {
    super();
    init(this, options, instance, create_fragment, safe_not_equal, {}, add_css);
  }
};
customElements.define("app-todo", create_custom_element(App, {}, [], [], true));
var app_default = App;
export {
  app_default as default
};

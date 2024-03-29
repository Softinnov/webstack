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
function toggle_class(element2, name, toggle) {
  element2.classList.toggle(name, !!toggle);
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
  append_styles(target, "svelte-5sbuja", 'h1.svelte-5sbuja{font-size:70px}p.svelte-5sbuja{font-size:large;left:45%;bottom:5%;position:fixed}.button.svelte-5sbuja{height:35px;width:35px;border:none;border-radius:5%;background-color:white;margin-right:3%}.button.svelte-5sbuja:hover{background-color:rgba(146, 146, 146, 0.381);cursor:pointer}.centered.svelte-5sbuja{width:25em;margin:auto;display:grid}.newTask.svelte-5sbuja{left:12%;margin-bottom:15%;margin-right:1%;position:relative}.todos.svelte-5sbuja{margin-top:10%}.ajout.svelte-5sbuja{height:35px;width:35px;border:none;border-radius:5%;background-color:white;margin-left:13%;position:relative}.ajout.svelte-5sbuja:hover{background-color:rgba(146, 146, 146, 0.381);cursor:pointer}.ajout.svelte-5sbuja:disabled{background-color:white;color:rgba(128, 128, 128, 0.836);cursor:default}.done.svelte-5sbuja{opacity:0.4}ul.svelte-5sbuja{max-height:15em;overflow:hidden;overflow-y:visible}li.svelte-5sbuja{display:flex;margin:3%}input[type="text"].svelte-5sbuja{flex:1;padding:0.5em;margin:-0.2em 0;border:none;font-size:large}');
}
function get_each_context(ctx, list, i) {
  const child_ctx = ctx.slice();
  child_ctx[11] = list[i];
  child_ctx[12] = list;
  child_ctx[13] = i;
  return child_ctx;
}
function create_each_block(ctx) {
  let li;
  let button0;
  let t1;
  let button1;
  let t3;
  let input;
  let t4;
  let mounted;
  let dispose;
  function input_input_handler_1() {
    ctx[8].call(
      input,
      /*each_value*/
      ctx[12],
      /*item_index*/
      ctx[13]
    );
  }
  return {
    c() {
      li = element("li");
      button0 = element("button");
      button0.textContent = "\u270F\uFE0F";
      t1 = space();
      button1 = element("button");
      button1.textContent = "\u{1F5D1}\uFE0F";
      t3 = space();
      input = element("input");
      t4 = space();
      attr(button0, "class", "button svelte-5sbuja");
      attr(button1, "class", "button svelte-5sbuja");
      attr(input, "type", "text");
      attr(input, "class", "svelte-5sbuja");
      attr(li, "class", "svelte-5sbuja");
      toggle_class(
        li,
        "done",
        /*item*/
        ctx[11].done
      );
    },
    m(target, anchor) {
      insert(target, li, anchor);
      append(li, button0);
      append(li, t1);
      append(li, button1);
      append(li, t3);
      append(li, input);
      set_input_value(
        input,
        /*item*/
        ctx[11].text
      );
      append(li, t4);
      if (!mounted) {
        dispose = [
          listen(button0, "click", function() {
            if (is_function(
              /*modify*/
              ctx[5](
                /*item*/
                ctx[11]
              )
            ))
              ctx[5](
                /*item*/
                ctx[11]
              ).apply(this, arguments);
          }),
          listen(button1, "click", function() {
            if (is_function(
              /*xclear*/
              ctx[4](
                /*item*/
                ctx[11]
              )
            ))
              ctx[4](
                /*item*/
                ctx[11]
              ).apply(this, arguments);
          }),
          listen(input, "input", input_input_handler_1),
          listen(input, "keydown", function() {
            if (is_function(
              /*handleKeydown*/
              ctx[6](
                /*item*/
                ctx[11],
                "modify"
              )
            ))
              ctx[6](
                /*item*/
                ctx[11],
                "modify"
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
      ctx[11].text) {
        set_input_value(
          input,
          /*item*/
          ctx[11].text
        );
      }
      if (dirty & /*todos*/
      1) {
        toggle_class(
          li,
          "done",
          /*item*/
          ctx[11].done
        );
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
  let div1;
  let h1;
  let t1;
  let div0;
  let input;
  let t2;
  let button;
  let t3;
  let button_disabled_value;
  let t4;
  let ul;
  let t5;
  let p;
  let t6;
  let t7;
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
      div1 = element("div");
      h1 = element("h1");
      h1.textContent = "My TodoList";
      t1 = space();
      div0 = element("div");
      input = element("input");
      t2 = space();
      button = element("button");
      t3 = text("\u2714\uFE0F");
      t4 = space();
      ul = element("ul");
      for (let i = 0; i < each_blocks.length; i += 1) {
        each_blocks[i].c();
      }
      t5 = space();
      p = element("p");
      t6 = text(
        /*remaining*/
        ctx[2]
      );
      t7 = text(" t\xE2ches restantes !");
      attr(h1, "class", "svelte-5sbuja");
      attr(input, "class", "newTask svelte-5sbuja");
      attr(input, "type", "text");
      attr(input, "placeholder", "Quoi d'autre?");
      attr(button, "class", "ajout svelte-5sbuja");
      button.disabled = button_disabled_value = /*nouvelleTache*/
      ctx[1] == "";
      attr(ul, "id", "todo-list");
      attr(ul, "class", "todos svelte-5sbuja");
      attr(p, "class", "svelte-5sbuja");
      attr(div1, "class", "centered svelte-5sbuja");
    },
    m(target, anchor) {
      insert(target, div1, anchor);
      append(div1, h1);
      append(div1, t1);
      append(div1, div0);
      append(div0, input);
      set_input_value(
        input,
        /*nouvelleTache*/
        ctx[1]
      );
      append(div0, t2);
      append(div0, button);
      append(button, t3);
      append(div1, t4);
      append(div1, ul);
      for (let i = 0; i < each_blocks.length; i += 1) {
        if (each_blocks[i]) {
          each_blocks[i].m(ul, null);
        }
      }
      append(div1, t5);
      append(div1, p);
      append(p, t6);
      append(p, t7);
      if (!mounted) {
        dispose = [
          listen(
            input,
            "input",
            /*input_input_handler*/
            ctx[7]
          ),
          listen(
            input,
            "keydown",
            /*handleKeydown*/
            ctx[6]
          ),
          listen(
            button,
            "click",
            /*add*/
            ctx[3]
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
      if (dirty & /*nouvelleTache*/
      2 && button_disabled_value !== (button_disabled_value = /*nouvelleTache*/
      ctx2[1] == "")) {
        button.disabled = button_disabled_value;
      }
      if (dirty & /*todos, handleKeydown, xclear, modify*/
      113) {
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
      4)
        set_data(
          t6,
          /*remaining*/
          ctx2[2]
        );
    },
    i: noop,
    o: noop,
    d(detaching) {
      if (detaching) {
        detach(div1);
      }
      destroy_each(each_blocks, detaching);
      mounted = false;
      run_all(dispose);
    }
  };
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
  function add() {
    let todo = { text: nouvelleTache };
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
      var url = `/${route}?text=${todo.text}`;
    } else {
      url = `/${route}?id=${todo.id}&text=${todo.text}`;
    }
    try {
      const reponse = await fetch(url, { method: "POST" });
      if (!reponse.ok) {
        const errorData = await reponse.text();
        alert(errorData);
        throw new Error(`Erreur lors de la requ\xEAte : ${reponse.status} ${reponse.statusText}`);
      }
      const result = await reponse.json();
      getTodoList();
    } catch (error) {
      console.error(`Une erreur s'est produite : ${error.message}`);
    }
  }
  function handleKeydown(e, item, action) {
    if (e.key == "Enter") {
      if (action == "modify") {
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
  function input_input_handler_1(each_value, item_index) {
    each_value[item_index].text = this.value;
    $$invalidate(0, todos);
  }
  $$self.$$.update = () => {
    if ($$self.$$.dirty & /*todos*/
    1) {
      $:
        $$invalidate(2, remaining = todos.length);
    }
  };
  return [
    todos,
    nouvelleTache,
    remaining,
    add,
    xclear,
    modify,
    handleKeydown,
    input_input_handler,
    input_input_handler_1
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

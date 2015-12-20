.. NanoGUI.go documentation master file, created by
   sphinx-quickstart on Sat Dec  5 23:10:08 2015.
   You can adapt this file completely to your liking, but it should at least
   contain the root `toctree` directive.

NanoGUI.go
======================================

.. image:: https://godoc.org/github.com/shibukawa/nanogui.go?status.svg
   :target: https://godoc.org/github.com/shibukawa/nanogui.go

.. raw:: html

   <iframe src="/nanogui.go/demo/" width="900" height="540"></iframe>

`Full Screen </nanogui.go/demo/>`_

NanoGUI.go is apure golang implementation of `NanoGUI <https://github.com/wjakob/nanogui>`_. NanoGUI is a minimalistic GUI library for OpenGL.

NanoGUI.go can run on desktops and browsers (I didn't test on other environemnt).

It uses the following components:

* `OpenGL/WebGL library <https://github.com/goxjs/gl>`_
* `Cross platform glfw wrapper <https://github.com/goxjs/glfw>`_
* `NanoVGo <https://shibukawa.github.io/nanovgo/>`_

Install
-------------

.. code-block:: bash

   $ go get github.com/shibukawa/nanogui.go

.. warning::

   Current code depends on unreleased version of GLFW to support IME (now sending PR for 3.2).

   Use the following packages before v3.2 would be released. If v3.2 will be released, **github.com/shibukawa/glfw\*** packages will be removed.

   - **github.com/go-gl/glfw/v3.1/glfw** → **gihtub.com/shibukawa/glfw-2/v3.2/glfw**
   - **github.com/goxjs/glfw** → **github.com/shibukawa/glfw**

API Reference
---------------

See `GoDoc <https://godoc.org/github.com/shibukawa/nanogui.go>`_

Porting Status
------------------

.. list-table::
   :widths: 15 15 10 10
   :header-rows: 1

   - * Category
     * Classes
     * Finished
     * Status
   - * Non-Visual types
     * Screen
     * ☑
     * except cursor feature
   - * 
     * Widget
     * ☑
     *
   - *
     * GroupLayout, BoxLayout, GridLayout, Advancedgridlayout
     * ☑
     *
   - *
     * Theme
     * ☑
     * Entype icons, Roboto fonts included
   - * Widgets
     * Window, Popup
     * ☑
     *
   - *
     * Label
     * ☑
     *
   - *
     * Button, ToolButton
     * ☑
     *
   - * 
     * PopupButton, ComboBox, ImagePanel, ColorPicker
     * ☑
     *
   - *
     * CheckBox
     * ☑
     *
   - *
     * Slider
     * ☑
     *
   - *
     * ProgressBar
     * ☑
     *
   - * 
     * TextBox, IntBox, FloatBox
     * ☑
     * It supports Emacs key bind like MacOS X, IME not supported.
   - *
     * ImageView, Graph
     * ☑
     *
   - *
     * VScrollPanel
     * ☑
     *
   - *
     * ColorWheel
     * ☑
     *
   - *
     * MessageDialog
     * ☐
     *
   - * Utitility
     * FormHelper
     * ☐
     *
   - * OpenGL Helper
     * GLShader, GLFramebuffer, Arcball
     * ☐
     *

Author
---------------

* `Yoshiki Shibukawa <https://github.com/shibukawa>`_

License
----------

NanoGUI.go is released under zlib license. But the following contents are under other licenses:

* Roboto fonts: Apatch 2 license
* Entypo icon: Creative Commons 4.0 BY-SA

Indices and tables
==================

* :ref:`genindex`
* :ref:`modindex`
* :ref:`search`


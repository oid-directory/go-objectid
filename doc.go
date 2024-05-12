/*
Package objectid implements ASN.1 Object Identifier types and methods.

# Features

  - Boundless [NumberForm] support
  - ASN.1 encoding and decoding of [DotNotation] instances -- without use of the [encoding/asn1] package
  - Flexible index support, allowing interrogation through negative indices without the risk of panic
  - Convenient Leaf, Parent and Root index alias methods, wherever applicable
  - Ge, Gt, Le, Lt, Equal comparison methods for interacting with [NumberForm] instances
  - Conversion friendly -- seamless [encoding/asn1.ObjectIdentifier] and [crypto/x509.OID] output support

# License

The go-objectid package is available under the terms of the MIT license.

For further details, see the LICENSE file within the root of the source repository.

# Boundless NumberForm

Instances of the [NumberForm] type are subject to no magnitude limits. This means that any given [NumberForm] may be set to any number, provided it is unsigned, regardless of its size.  In particular, this is necessary to support ITU-T Rec. X.667 OIDs, which are UUID based and require 128-bit unsigned integer support.

Be aware that the [NumberForm] type is based upon [math/big.Int]. Therefore, given sufficiently large values, performance penalties may appear during routine operations. Keep in mind that most OIDs do not bear such large values.

# ASN.1 Codec

The ASN.1 encoding and decoding scheme implemented in this package does not use the [encoding/asn1] package, nor does it utilize any 3rd party ASN.1 implementation. It is a custom, bidirectional implementation.

ASN.1 codec features are written solely for [DotNotation] related cases (e.g.: "1.3.6.1.4.1.56521"), and allow the following:

  - Encoding (marshaling) of a [DotNotation] into an ASN.1 encoded value ([]byte{...})
  - Decoding (unmarshaling) of encoded values into an unpopulated [DotNotation] instance

Encoding of non-minimal values -- such as root arcs "0", "1" and "2" alone -- is not supported.  Some ASN.1 implementations precariously treat certain OIDs, such as "0" and "0.0" the same, likely for support reasons. This results in ambiguity when handling pre-encoded bytes in an obverse scenario, and is in violation of ITU-T Rec. X.690 regarding the proper encoding of an ASN.1 OBJECT IDENTIFIER.

In short, codec functions will only operate successfully when given [DotNotation] comprised of two (2) or more [NumberForm] instances.
*/
package objectid

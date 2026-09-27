#!/usr/bin/perl
# Dump ExifTool tag tables to simplified YAML for imagemeta meta/exif/tags.
# Source: ExifTool perl modules (installed version, see exiftool_version below).
# Output: one YAML file per table, Go-codegen friendly (no CODE refs dumped).
use strict;
use warnings;

my $LIB = $ENV{EXIFTOOL_LIB} || '/opt/homebrew/Cellar/exiftool/13.55/libexec/lib/perl5';
push @INC, $LIB;

use Image::ExifTool::Exif;
use Image::ExifTool::GPS;
use Image::ExifTool::Canon;
use Image::ExifTool::Nikon;
use Image::ExifTool::Sony;
use Image::ExifTool::Panasonic;
use Image::ExifTool::Apple;

my $OUT = $ENV{TAGS_OUT} || 'meta/exif/tags';

sub yaml_escape {
    my ($s) = @_;
    return '""' unless defined $s && length $s;
    $s =~ s/\\/\\\\/g;
    $s =~ s/"/\\"/g;
    $s =~ s/\n/\\n/g;
    $s =~ s/\r/\\r/g;
    $s =~ s/\t/\\t/g;
    # trim to keep YAML readable
    if (length $s > 220) { $s = substr($s, 0, 217) . '...'; }
    return '"' . $s . '"';
}

sub short_conv {
    my ($v) = @_;
    return 'none' unless defined $v;
    my $r = ref $v;
    return 'code' if $r eq 'CODE';
    return 'hash' if $r eq 'HASH';
    return 'array' if $r eq 'ARRAY';
    # scalar: truncate, classify
    my $s = "$v";
    $s =~ s/^\s+|\s+$//g;
    return 'none' if $s eq '';
    if ($s =~ /^sub\s*\{/ || $s =~ /^\$val/ || $s =~ /Image::ExifTool/) { return 'code'; }
    if (length $s > 120) { $s = substr($s, 0, 117) . '...'; }
    return "scalar:" . $s;
}

sub is_tag_key {
    my ($k) = @_;
    return 1 if $k =~ /^(0x[0-9a-fA-F]+|\d+)$/;
    return 0;
}

sub key_to_id {
    my ($k) = @_;
    if ($k =~ /^0x/i) { return sprintf("0x%04x", hex($k)); }
    return sprintf("0x%04x", $k);
}

sub dump_table {
    my (%a) = @_;
    my ($file, $version, $module, $pm, $table, $href, $docs, $groups) = @a{qw(file version module pm table href docs groups)};
    open my $fh, '>', $file or die "open $file: $!";
    print $fh "# Code generated from ExifTool. DO NOT EDIT.\n";
    print $fh "# Source: $module $table ($pm), ExifTool $version\n";
    print $fh "# Docs: $docs\n";
    print $fh "# Generator: meta/exif/tags/dump_exiftool_tables.pl (see README.md)\n";
    print $fh "exiftool_version: \"$version\"\n";
    print $fh "source:\n";
    print $fh "  module: \"$module\"\n";
    print $fh "  table: \"$table\"\n";
    print $fh "  file: \"$pm\"\n";
    print $fh "  docs: \"$docs\"\n";
    if ($groups && %$groups) {
        print $fh "groups:\n";
        for my $g (sort { $a <=> $b } keys %$groups) {
            print $fh "  $g: " . yaml_escape($groups->{$g}) . "\n";
        }
    }
    my @keys = sort { hex($a) <=> hex($b) } grep { is_tag_key($_) } keys %$href;
    print $fh "source_tag_count: " . scalar(@keys) . "\n";
    my @written;
    for my $k (@keys) {
        my $e0 = $href->{$k};
        next unless ref $e0 eq 'HASH';
        my $n0 = $e0->{Name} // '';
        push @written, $k if $n0 =~ /^[A-Za-z]/;
    }
    print $fh "tag_count: " . scalar(@written) . "\n";
    print $fh "tags:\n";
    for my $k (@keys) {
        my $e = $href->{$k};
        next unless ref $e eq 'HASH';
        my $id = key_to_id($k);
        my $name = $e->{Name} // '';
        # skip binary-data / subdirectory-only aliases without a name
        next unless $name =~ /^[A-Za-z]/;
        my $writable = $e->{Writable} // '';
        my $wgroup = $e->{WriteGroup} // $e->{Group} // '';
        my $count = exists $e->{Count} ? $e->{Count} : '';
        my $desc = $e->{Description} // $e->{Notes} // '';
        $desc =~ s/\s+/ /g;
        my $subdir = exists $e->{SubDirectory} ? 'true' : 'false';
        my $subdir_table = '';
        if (ref $e->{SubDirectory} eq 'HASH' && $e->{SubDirectory}{TagTable}) {
            $subdir_table = $e->{SubDirectory}{TagTable};
            $subdir_table =~ s/^Image::ExifTool:://;
        }
        my $pconv = short_conv($e->{PrintConv});
        my $vconv = short_conv($e->{ValueConv});
        print $fh "  - id: $id\n";
        print $fh "    name: " . yaml_escape($name) . "\n";
        print $fh "    writable: " . yaml_escape("$writable") . "\n";
        print $fh "    group: " . yaml_escape("$wgroup") . "\n";
        print $fh "    count: " . yaml_escape("$count") . "\n";
        print $fh "    description: " . yaml_escape($desc) . "\n";
        print $fh "    print_conv: " . yaml_escape($pconv) . "\n";
        print $fh "    value_conv: " . yaml_escape($vconv) . "\n";
        print $fh "    subdirectory: $subdir\n";
        if ($subdir_table) {
            print $fh "    subdirectory_table: " . yaml_escape($subdir_table) . "\n";
        }
    }
    close $fh;
    print "wrote $file\n";
}

my $ver = $Image::ExifTool::VERSION // 'unknown';

dump_table(
    file => "$OUT/exif.yaml", version => $ver,
    module => 'Image::ExifTool::Exif', pm => 'Exif.pm', table => '%Main',
    href => \%Image::ExifTool::Exif::Main,
    docs => 'https://exiftool.org/TagNames/EXIF.html',
    groups => $Image::ExifTool::Exif::Main{GROUPS},
);
dump_table(
    file => "$OUT/gps.yaml", version => $ver,
    module => 'Image::ExifTool::GPS', pm => 'GPS.pm', table => '%Main',
    href => \%Image::ExifTool::GPS::Main,
    docs => 'https://exiftool.org/TagNames/GPS.html',
    groups => $Image::ExifTool::GPS::Main{GROUPS},
);
dump_table(
    file => "$OUT/canon.yaml", version => $ver,
    module => 'Image::ExifTool::Canon', pm => 'Canon.pm', table => '%Main',
    href => \%Image::ExifTool::Canon::Main,
    docs => 'https://exiftool.org/TagNames/Canon.html',
    groups => $Image::ExifTool::Canon::Main{GROUPS},
);
dump_table(
    file => "$OUT/nikon.yaml", version => $ver,
    module => 'Image::ExifTool::Nikon', pm => 'Nikon.pm', table => '%Main',
    href => \%Image::ExifTool::Nikon::Main,
    docs => 'https://exiftool.org/TagNames/Nikon.html',
    groups => $Image::ExifTool::Nikon::Main{GROUPS},
);
dump_table(
    file => "$OUT/sony.yaml", version => $ver,
    module => 'Image::ExifTool::Sony', pm => 'Sony.pm', table => '%Main',
    href => \%Image::ExifTool::Sony::Main,
    docs => 'https://exiftool.org/TagNames/Sony.html',
    groups => $Image::ExifTool::Sony::Main{GROUPS},
);
dump_table(
    file => "$OUT/panasonic.yaml", version => $ver,
    module => 'Image::ExifTool::Panasonic', pm => 'Panasonic.pm', table => '%Main',
    href => \%Image::ExifTool::Panasonic::Main,
    docs => 'https://exiftool.org/TagNames/Panasonic.html',
    groups => $Image::ExifTool::Panasonic::Main{GROUPS},
);
dump_table(
    file => "$OUT/apple.yaml", version => $ver,
    module => 'Image::ExifTool::Apple', pm => 'Apple.pm', table => '%Main',
    href => \%Image::ExifTool::Apple::Main,
    docs => 'https://exiftool.org/TagNames/Apple.html',
    groups => $Image::ExifTool::Apple::Main{GROUPS},
);

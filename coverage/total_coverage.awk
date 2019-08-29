{
    if ($0 ~ /^total\:/)
        print "coverage: " $3 " of statements";
}
